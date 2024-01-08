package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"goa.design/clue/clue"
	"goa.design/clue/debug"
	"goa.design/clue/health"
	"goa.design/clue/log"
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
		oteladdr             = flag.String("otel-addr", ":4317", "OpenTelemetry collector listen address")
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
	ctx := log.Context(context.Background(), log.WithFormat(format), log.WithFunc(log.Span))
	if *debugf {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}

	// 2. Setup isntrumentation
	spanExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(*oteladdr),
		otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(ctx, err, "failed to initialize tracing")
	}
	defer func() {
		// Create new context in case the parent context has been canceled.
		ctx := log.Context(context.Background(), log.WithFormat(format))
		err := spanExporter.Shutdown(ctx)
		if err != nil {
			log.Errorf(ctx, err, "failed to shutdown tracing")
		}
	}()
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(*oteladdr),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(ctx, err, "failed to initialize metrics")
	}
	defer func() {
		// Create new context in case the parent context has been canceled.
		ctx := log.Context(context.Background(), log.WithFormat(format))
		err := metricExporter.Shutdown(ctx)
		if err != nil {
			log.Errorf(ctx, err, "failed to shutdown metrics")
		}
	}()
	cfg, err := clue.NewConfig(ctx,
		gentester.ServiceName,
		gentester.APIVersion,
		metricExporter,
		spanExporter,
	)
	if err != nil {
		log.Fatalf(ctx, err, "failed to initialize instrumentation")
	}
	clue.ConfigureOpenTelemetry(ctx, cfg)

	// 3. Create clients
	lcc, err := grpc.DialContext(ctx,
		*locatorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(log.UnaryClientInterceptor()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Fatalf(ctx, err, "failed to connect to locator")
	}
	lc := locator.New(lcc)
	fcc, err := grpc.DialContext(ctx, *forecasterAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(log.UnaryClientInterceptor()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Fatalf(ctx, err, "failed to connect to forecast")
	}
	fc := forecaster.New(fcc)

	// 4. Create service & endpoints
	svc := tester.New(lc, fc)
	endpoints := gentester.NewEndpoints(svc)
	endpoints.Use(debug.LogPayloads())
	endpoints.Use(log.Endpoint)

	// 5. Create transport
	server := gengrpc.New(endpoints, nil)
	grpcsvr := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			log.UnaryServerInterceptor(ctx),
			debug.UnaryServerInterceptor()),
		grpc.StatsHandler(otelgrpc.NewServerHandler()))
	genpb.RegisterTesterServer(grpcsvr, server)
	reflection.Register(grpcsvr)
	for svc, info := range grpcsvr.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Print(ctx, log.KV{K: "method", V: svc + "/" + m.Name})
		}
	}

	// 6. Setup health check, metrics and debug endpoints
	mux := http.NewServeMux()
	debug.MountDebugLogEnabler(mux)
	debug.MountPprofHandlers(mux)
	check := health.Handler(health.NewChecker(
		health.NewPinger("locator", *locatorHealthAddr),
		health.NewPinger("forecaster", *forecasterHealthAddr)))
	mux.Handle("/healthz", check)
	mux.Handle("/livez", check)
	httpsvr := &http.Server{Addr: *httpaddr, Handler: mux}

	// 7. Start gRPC and HTTP servers
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("signal: %s", <-c)
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
	if err := <-errc; err != nil && !strings.HasPrefix(err.Error(), "signal:") {
		log.Errorf(ctx, err, "exiting")
	}
	cancel()
	wg.Wait()
	log.Printf(ctx, "exited")
}
