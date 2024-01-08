package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"goa.design/clue/clue"
	"goa.design/clue/debug"
	"goa.design/clue/health"
	"goa.design/clue/log"
	goahttp "goa.design/goa/v3/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"goa.design/clue/example/weather/services/front"
	"goa.design/clue/example/weather/services/front/clients/forecaster"
	"goa.design/clue/example/weather/services/front/clients/locator"
	"goa.design/clue/example/weather/services/front/clients/tester"
	genfront "goa.design/clue/example/weather/services/front/gen/front"
	genhttp "goa.design/clue/example/weather/services/front/gen/http/front/server"
)

func main() {
	var (
		httpListenAddr       = flag.String("http-addr", ":8084", "HTTP listen address")
		metricsListenAddr    = flag.String("metrics-addr", ":8085", "metrics listen address")
		forecasterAddr       = flag.String("forecaster-addr", ":8080", "Forecaster service address")
		forecasterHealthAddr = flag.String("forecaster-health-addr", ":8081", "Forecaster service health-check address")
		locatorAddr          = flag.String("locator-addr", ":8082", "Locator service address")
		locatorHealthAddr    = flag.String("locator-health-addr", ":8083", "Locator service health-check address")
		coladdr              = flag.String("otel-addr", ":4317", "OpenTelemtry collector listen address")
		debugf               = flag.Bool("debug", false, "Enable debug logs")
		testerAddr           = flag.String("tester-addr", ":8090", "Tester service address")
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

	// 2. Setup instrumentation
	spanExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(*coladdr),
		otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(ctx, err, "failed to initialize tracing")
	}
	defer func() {
		// Create new context in case the parent context has been canceled.
		ctx := log.Context(context.Background(), log.WithFormat(format))
		if err := spanExporter.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown tracing")
		}
	}()
	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(*coladdr),
		otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(ctx, err, "failed to initialize metrics")
	}
	defer func() {
		// Create new context in case the parent context has been canceled.
		ctx := log.Context(context.Background(), log.WithFormat(format))
		if err := metricExporter.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown metrics")
		}
	}()
	cfg, err := clue.NewConfig(ctx,
		genfront.ServiceName,
		genfront.APIVersion,
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
	tcc, err := grpc.DialContext(ctx, *testerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(log.UnaryClientInterceptor()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Errorf(ctx, err, "failed to connect to tester")
		os.Exit(1)
	}
	tc := tester.New(tcc)

	// 4. Create service & endpoints
	svc := front.New(fc, lc, tc)
	endpoints := genfront.NewEndpoints(svc)
	endpoints.Use(debug.LogPayloads())
	endpoints.Use(log.Endpoint)

	// 5. Create transport
	mux := goahttp.NewMuxer()
	debug.MountDebugLogEnabler(debug.Adapt(mux))
	debug.MountPprofHandlers(debug.Adapt(mux))
	handler := otelhttp.NewHandler(mux, genfront.ServiceName) // 3. Add OpenTelemetry instrumentation
	handler = debug.HTTP()(handler)                           // 2. Add debug endpoints
	handler = log.HTTP(ctx)(handler)                          // 1. Add logger to request context
	server := genhttp.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	genhttp.Mount(mux, server)
	for _, m := range server.Mounts {
		log.Print(ctx, log.KV{K: "method", V: m.Method}, log.KV{K: "endpoint", V: m.Verb + " " + m.Pattern})
	}
	httpServer := &http.Server{Addr: *httpListenAddr, Handler: handler}

	// 6. Mount health check & metrics on separate HTTP server (different listen port)
	// No testerHealthAddr pinger because we don't want the whole system to die just because
	// tester isn't healthy for some reason
	check := health.Handler(health.NewChecker(
		health.NewPinger("locator", *locatorHealthAddr),
		health.NewPinger("forecaster", *forecasterHealthAddr)))
	check = log.HTTP(ctx)(check).(http.HandlerFunc) // Log health-check errors
	http.Handle("/healthz", check)
	http.Handle("/livez", check)
	metricsServer := &http.Server{Addr: *metricsListenAddr}

	// 7. Start HTTP server
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
			log.Printf(ctx, "HTTP server listening on %s", *httpListenAddr)
			errc <- httpServer.ListenAndServe()
		}()

		go func() {
			log.Printf(ctx, "Metrics server listening on %s", *metricsListenAddr)
			errc <- metricsServer.ListenAndServe()
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down HTTP servers")

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown HTTP server")
		}
		if err := metricsServer.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown metrics server")
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
