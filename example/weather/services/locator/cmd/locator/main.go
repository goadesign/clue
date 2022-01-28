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

	"github.com/goadesign/clue/health"
	"github.com/goadesign/clue/instrument"
	"github.com/goadesign/clue/log"
	"github.com/goadesign/clue/trace"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	goagrpcmiddleware "goa.design/goa/v3/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/goadesign/clue/example/weather/services/locator"
	"github.com/goadesign/clue/example/weather/services/locator/clients/ipapi"
	genpb "github.com/goadesign/clue/example/weather/services/locator/gen/grpc/locator/pb"
	gengrpc "github.com/goadesign/clue/example/weather/services/locator/gen/grpc/locator/server"
	genlocator "github.com/goadesign/clue/example/weather/services/locator/gen/locator"
)

func main() {
	var (
		grpcaddr  = flag.String("grpc-addr", ":8082", "gRPC listen address")
		httpaddr  = flag.String("http-addr", ":8083", "HTTP listen address (health checks)")
		agentaddr = flag.String("agent-addr", ":4317", "Grafana agent listen address")
		debug     = flag.Bool("debug", false, "Enable debug logs")
	)
	flag.Parse()

	// 1. Create logger
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format))
	ctx = log.With(ctx, "svc", genlocator.ServiceName)
	if *debug {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debug(ctx, "debug logs enabled")
	}

	// 2. Setup tracing
	log.Debug(ctx, "connecting to Grafana agent...", "addr", *agentaddr)
	conn, err := grpc.DialContext(ctx, *agentaddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.Error(ctx, "failed to connect to Grafana agent", "err", err)
		os.Exit(1)
	}
	log.Debug(ctx, "connected to Grafana agent", "addr", *agentaddr)
	ctx, err = trace.Context(ctx, genlocator.ServiceName, conn)
	if err != nil {
		log.Error(ctx, "failed to initialize tracing", "err", err)
		os.Exit(1)
	}

	// 3. Setup instrumentation
	ctx = instrument.Context(ctx, genlocator.ServiceName)

	// 4. Create clients
	c := &http.Client{Transport: trace.Client(ctx, http.DefaultTransport)}
	ipc := ipapi.New(c)

	// 5. Create service & endpoints
	svc := locator.New(ipc)
	endpoints := genlocator.NewEndpoints(svc)

	// 6. Create transport
	server := gengrpc.New(endpoints, nil)
	grpcsvr := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			goagrpcmiddleware.UnaryRequestID(),
			log.UnaryServerInterceptor(ctx),
			goagrpcmiddleware.UnaryServerLog(log.Adapt(ctx)),
			trace.UnaryServerInterceptor(ctx),
			instrument.UnaryServerInterceptor(ctx),
		))
	genpb.RegisterLocatorServer(grpcsvr, server)
	reflection.Register(grpcsvr)
	for svc, info := range grpcsvr.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Print(ctx, "mount", "method", svc+"/"+m.Name)
		}
	}

	// 7. Start health check
	check := log.HTTP(ctx)(health.Handler(health.NewChecker(ipc)))
	http.Handle("/healthz", check)
	http.Handle("/livez", check)
	http.Handle("/metrics", instrument.Handler(ctx))
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
			log.Print(ctx, "HTTP server listening", "addr", httpsvr.Addr)
			errc <- httpsvr.ListenAndServe()
		}()

		var l net.Listener
		go func() {
			var err error
			l, err = net.Listen("tcp", *grpcaddr)
			if err != nil {
				errc <- err
			}
			log.Print(ctx, "gRPC server listening", "addr", l.Addr())
			errc <- grpcsvr.Serve(l)
		}()

		<-ctx.Done()
		log.Print(ctx, "shutting down HTTP and gRPC servers")

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		grpcsvr.GracefulStop()
		httpsvr.Shutdown(ctx)
	}()

	// Cleanup
	log.Print(ctx, "exiting", "err", <-errc)
	cancel()
	wg.Wait()
	log.Print(ctx, "exited")
}
