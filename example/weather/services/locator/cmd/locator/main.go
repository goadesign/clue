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

	"github.com/crossnokaye/micro/health"
	"github.com/crossnokaye/micro/instrument"
	"github.com/crossnokaye/micro/log"
	"github.com/crossnokaye/micro/trace"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	goagrpcmiddleware "goa.design/goa/v3/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/crossnokaye/micro/example/weather/services/locator"
	"github.com/crossnokaye/micro/example/weather/services/locator/clients/ipapi"
	genpb "github.com/crossnokaye/micro/example/weather/services/locator/gen/grpc/locator/pb"
	gengrpc "github.com/crossnokaye/micro/example/weather/services/locator/gen/grpc/locator/server"
	genlocator "github.com/crossnokaye/micro/example/weather/services/locator/gen/locator"
)

func main() {
	var (
		grpcListenAddr = flag.String("grpc", ":8082", "gRPC listen address")
		httpListenAddr = flag.String("http", ":8083", "HTTP listen address (health checks)")
		collectorAddr  = flag.String("coladdr", ":4317", "OpenTelemetry remote collector address")
		debug          = flag.Bool("debug", false, "Enable debug logs")
	)
	flag.Parse()

	// 1. Create logger
	ctx := log.With(log.Context(context.Background()), "svc", genlocator.ServiceName)
	if *debug {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debug(ctx, "debug logs enabled")
	}

	// 2. Setup tracing
	log.Debug(ctx, "connecting to OpenTelemetry collector...", "addr", *collectorAddr)
	conn, err := grpc.DialContext(ctx, *collectorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.Error(ctx, "failed to connect to OpenTelementry collector", "err", err)
		os.Exit(1)
	}
	log.Debug(ctx, "connected to OpenTelemetry collector", "addr", *collectorAddr)
	ctx, err = trace.Context(ctx, genlocator.ServiceName, conn)
	if err != nil {
		log.Error(ctx, "failed to initialize tracing", "err", err)
		os.Exit(1)
	}

	// 3. Create clients
	c := &http.Client{Transport: trace.Client(ctx, http.DefaultTransport)}
	ipc := ipapi.New(c)

	// 4. Create service & endpoints
	svc := locator.New(ipc)
	endpoints := genlocator.NewEndpoints(svc)
	endpoints.Use(log.Init(ctx))

	// 5. Create transport
	server := gengrpc.New(endpoints, nil)
	grpcsvr := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			goagrpcmiddleware.UnaryRequestID(),
			goagrpcmiddleware.UnaryServerLog(log.Adapt(ctx)),
			trace.UnaryServerInterceptor(ctx),
			instrument.UnaryServerInterceptor(ctx, genlocator.ServiceName),
		))
	genpb.RegisterLocatorServer(grpcsvr, server)
	reflection.Register(grpcsvr)
	for svc, info := range grpcsvr.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Print(ctx, "mount", "method", svc+"/"+m.Name)
		}
	}

	// 6. Start health check
	check := log.HTTP(ctx)(health.Handler(health.NewChecker(ipc)))
	http.Handle("/healthz", check)
	http.Handle("/livez", check)
	http.Handle("/metrics", instrument.Handler(ctx))
	httpsvr := &http.Server{Addr: *httpListenAddr}

	// 7. Start gRPC and HTTP servers
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
			log.Print(ctx, "HTTP server listening", "addr", *httpListenAddr)
			errc <- httpsvr.ListenAndServe()
		}()

		var l net.Listener
		go func() {
			var err error
			l, err = net.Listen("tcp", *grpcListenAddr)
			if err != nil {
				errc <- err
			}
			log.Print(ctx, "gRPC server listening", "addr", *grpcListenAddr)
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
