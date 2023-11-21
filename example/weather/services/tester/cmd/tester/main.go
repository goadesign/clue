package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"goa.design/clue/debug"
	"goa.design/clue/health"
	"goa.design/clue/log"
	"goa.design/clue/metrics"
	"goa.design/clue/trace"
	goagrpcmiddleware "goa.design/goa/v3/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"goa.design/clue/example/weather/services/tester"
	"goa.design/clue/example/weather/services/tester/clients/forecaster"
	"goa.design/clue/example/weather/services/tester/clients/locator"
	genpb "goa.design/clue/example/weather/services/tester/gen/grpc/tester/pb"
	gengrpc "goa.design/clue/example/weather/services/tester/gen/grpc/tester/server"
	gentester "goa.design/clue/example/weather/services/tester/gen/tester"
)

func main() {
	var (
		grpcaddr             = flag.String("grpc-addr", ":8090", "gRPC listen address")
		httpaddr             = flag.String("http-addr", ":8091", "HTTP listen address (health checks and metrics)")
		agentaddr            = flag.String("agent-addr", ":4317", "Grafana agent listen address")
		forecasterAddr       = flag.String("forecaster-addr", ":8080", "Forecaster service address")
		forecasterHealthAddr = flag.String("forecaster-health-addr", ":8081", "Forecaster service health-check address")
		locatorAddr          = flag.String("locator-addr", ":8082", "Locator service address")
		locatorHealthAddr    = flag.String("locator-health-addr", ":8083", "Locator service health-check address")
		debugf               = flag.Bool("debug", false, "Enable debug logs")
	)
	flag.Parse()

	// 1. Create logger
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format), log.WithFunc(trace.Log))
	ctx = log.With(ctx, log.KV{K: "svc", V: gentester.ServiceName})
	if *debugf {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}

	// 2. Setup tracing
	log.Debugf(ctx, "connecting to Grafana agent %s", *agentaddr)
	conn, err := grpc.DialContext(ctx, *agentaddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.Errorf(ctx, err, "failed to connect to Grafana agent")
		os.Exit(1)
	}
	log.Debugf(ctx, "connected to Grafana agent %s", *agentaddr)
	ctx, err = trace.Context(ctx, gentester.ServiceName, trace.WithGRPCExporter(conn))
	if err != nil {
		log.Errorf(ctx, err, "failed to initialize tracing")
		os.Exit(1)
	}

	// 3. Setup metrics
	ctx = metrics.Context(ctx, gentester.ServiceName)

	// 4. Create clients
	lcc, err := grpc.DialContext(ctx, *locatorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			trace.UnaryClientInterceptor(ctx),
			log.UnaryClientInterceptor()))
	if err != nil {
		log.Errorf(ctx, err, "failed to connect to locator")
		os.Exit(1)
	}
	lc := locator.New(lcc)
	fcc, err := grpc.DialContext(ctx, *forecasterAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			trace.UnaryClientInterceptor(ctx),
			log.UnaryClientInterceptor()))
	if err != nil {
		log.Errorf(ctx, err, "failed to connect to forecast")
		os.Exit(1)
	}
	fc := forecaster.New(fcc)

	// 5. Create service & endpoints
	svc := tester.New(lc, fc)
	endpoints := gentester.NewEndpoints(svc)
	endpoints.Use(debug.LogPayloads())
	endpoints.Use(log.Endpoint)

	// 6. Create transport
	server := gengrpc.New(endpoints, nil)
	grpcsvr := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			goagrpcmiddleware.UnaryRequestID(),
			log.UnaryServerInterceptor(ctx),
			debug.UnaryServerInterceptor(),
			goagrpcmiddleware.UnaryServerLogContext(log.AsGoaMiddlewareLogger),
			trace.UnaryServerInterceptor(ctx),
			metrics.UnaryServerInterceptor(ctx),
		))
	genpb.RegisterTesterServer(grpcsvr, server)
	reflection.Register(grpcsvr)
	for svc, info := range grpcsvr.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Print(ctx, log.KV{K: "method", V: svc + "/" + m.Name})
		}
	}

	// 7. Setup health check, metrics and debug endpoints
	check := health.Handler(health.NewChecker(
		health.NewPinger("locator", *locatorHealthAddr),
		health.NewPinger("forecaster", *forecasterHealthAddr)))
	mux := http.NewServeMux()
	debug.MountDebugLogEnabler(mux)
	debug.MountPprofHandlers(mux)
	mux.Handle("/healthz", check)
	mux.Handle("/livez", check)
	mux.Handle("/metrics", metrics.Handler(ctx))
	httpsvr := &http.Server{Addr: *httpaddr, Handler: mux}

	// 8. Start gRPC and HTTP servers
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		go func() {
			log.Printf(ctx, "HTTP server listening on %s", httpsvr.Addr)
			errc <- httpsvr.ListenAndServe()
		}()

		var l net.Listener
		go func() {
			var err error
			l, err = net.Listen("tcp", *grpcaddr)
			if err != nil {
				errc <- err
			}
			log.Printf(ctx, "gRPC server listening on %s", l.Addr())
			errc <- grpcsvr.Serve(l)
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down HTTP and gRPC servers")

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		grpcsvr.GracefulStop()
		if err := httpsvr.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown HTTP server")
		}
	}()

	// Cleanup
	if err := <-errc; err != nil {
		log.Errorf(ctx, err, "exiting")
	}
	cancel()
	wg.Wait()
	log.Printf(ctx, "exited")
}
