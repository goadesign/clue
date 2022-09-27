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

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"goa.design/clue/health"
	"goa.design/clue/log"
	"goa.design/clue/metrics"
	"goa.design/clue/trace"
	goagrpcmiddleware "goa.design/goa/v3/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"goa.design/clue/example/weather/services/forecaster"
	"goa.design/clue/example/weather/services/forecaster/clients/weathergov"
	genforecaster "goa.design/clue/example/weather/services/forecaster/gen/forecaster"
	genpb "goa.design/clue/example/weather/services/forecaster/gen/grpc/forecaster/pb"
	gengrpc "goa.design/clue/example/weather/services/forecaster/gen/grpc/forecaster/server"
)

func main() {
	var (
		grpcaddr  = flag.String("grpc-addr", ":8080", "gRPC listen address")
		httpaddr  = flag.String("http-addr", ":8081", "HTTP listen address (health checks and metrics)")
		agentaddr = flag.String("agent-addr", ":4317", "Grafana agent listen address")
		debug     = flag.Bool("debug", false, "Enable debug logs")
	)
	flag.Parse()

	// 1. Create logger
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format), log.WithFunc(trace.Log))
	ctx = log.With(ctx, log.KV{"svc", genforecaster.ServiceName})
	if *debug {
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
	ctx, err = trace.Context(ctx, genforecaster.ServiceName, trace.WithGRPCExporter(conn))
	if err != nil {
		log.Errorf(ctx, err, "failed to initialize tracing")
		os.Exit(1)
	}

	// 3. Setup metrics
	ctx = metrics.Context(ctx, genforecaster.ServiceName)

	// 4. Create clients
	c := &http.Client{Transport: trace.Client(ctx, http.DefaultTransport)}
	wc := weathergov.New(log.Client(c))

	// 5. Create service & endpoints
	svc := forecaster.New(wc)
	endpoints := genforecaster.NewEndpoints(svc)

	// 6. Create transport
	server := gengrpc.New(endpoints, nil)
	grpcsvr := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			goagrpcmiddleware.UnaryRequestID(),
			log.UnaryServerInterceptor(ctx),
			goagrpcmiddleware.UnaryServerLogContext(log.AsGoaMiddlewareLogger),
			metrics.UnaryServerInterceptor(ctx),
			trace.UnaryServerInterceptor(ctx),
		))
	genpb.RegisterForecasterServer(grpcsvr, server)
	reflection.Register(grpcsvr)
	for svc, info := range grpcsvr.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Print(ctx, log.KV{"method", svc + "/" + m.Name})
		}
	}

	// 7. Setup health check and metrics
	check := log.HTTP(ctx)(health.Handler(health.NewChecker(wc)))
	http.Handle("/healthz", check)
	http.Handle("/livez", check)
	http.Handle("/metrics", metrics.Handler(ctx))
	httpsvr := &http.Server{Addr: *httpaddr}

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
			log.Printf(ctx, "HTTP server listening on %s", *httpaddr)
			errc <- httpsvr.ListenAndServe()
		}()

		var l net.Listener
		go func() {
			var err error
			l, err = net.Listen("tcp", *grpcaddr)
			if err != nil {
				errc <- err
			}
			log.Printf(ctx, "gRPC server listening on %s", *grpcaddr)
			errc <- grpcsvr.Serve(l)
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down HTTP and gRPC servers")

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		grpcsvr.GracefulStop()
		httpsvr.Shutdown(ctx)
	}()

	// Cleanup
	if err := <-errc; err != nil {
		log.Errorf(ctx, err, "exiting")
	}
	cancel()
	wg.Wait()
	log.Printf(ctx, "exited")
}
