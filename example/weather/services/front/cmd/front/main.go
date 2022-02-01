package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"goa.design/clue/health"
	"goa.design/clue/log"
	"goa.design/clue/metrics"
	"goa.design/clue/trace"
	goahttp "goa.design/goa/v3/http"
	goahttpmiddleware "goa.design/goa/v3/http/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"goa.design/clue/example/weather/services/front"
	"goa.design/clue/example/weather/services/front/clients/forecaster"
	"goa.design/clue/example/weather/services/front/clients/locator"
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
		agentaddr            = flag.String("agent-addr", ":4317", "Grafana agent listen address")
		debug                = flag.Bool("debug", false, "Enable debug logs")
	)
	flag.Parse()

	// 1. Create logger
	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format))
	ctx = log.With(ctx, "svc", genfront.ServiceName)
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
	ctx, err = trace.Context(ctx, genfront.ServiceName, conn)
	if err != nil {
		log.Error(ctx, "failed to initialize tracing", "err", err)
		os.Exit(1)
	}

	// 3. Setup metrics
	ctx = metrics.Context(ctx, genfront.ServiceName)

	// 3. Create clients
	lcc, err := grpc.DialContext(ctx, *locatorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(trace.UnaryClientInterceptor(ctx)))
	if err != nil {
		log.Error(ctx, "failed to connect to locator", "err", err)
		os.Exit(1)
	}
	lc := locator.New(lcc)
	fcc, err := grpc.DialContext(ctx, *forecasterAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(trace.UnaryClientInterceptor(ctx)))
	if err != nil {
		log.Error(ctx, "failed to connect to forecast", "err", err)
		os.Exit(1)
	}
	fc := forecaster.New(fcc)

	// 4. Create service & endpoints
	svc := front.New(fc, lc)
	endpoints := genfront.NewEndpoints(svc)

	// 5. Create transport
	mux := goahttp.NewMuxer()
	server := genhttp.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	genhttp.Mount(mux, server)
	handler := trace.HTTP(ctx)(mux)
	handler = metrics.HTTP(ctx)(handler)
	handler = goahttpmiddleware.Log(log.Adapt(ctx))(handler)
	handler = log.HTTP(ctx)(handler)
	handler = goahttpmiddleware.RequestID()(handler)
	for _, m := range server.Mounts {
		log.Print(ctx, "mount", "method", m.Method, "verb", m.Verb, "path", m.Pattern)
	}
	httpServer := &http.Server{Addr: *httpListenAddr, Handler: handler}

	// 6. Mount health check & metrics on separate HTTP server (different listen port)
	check := health.Handler(health.NewChecker(
		health.NewPinger("locator", "http", *locatorHealthAddr),
		health.NewPinger("forecaster", "http", *forecasterHealthAddr)))
	check = log.HTTP(ctx)(check).(http.HandlerFunc)
	http.Handle("/healthz", check)
	http.Handle("/livez", check)
	http.Handle("/metrics", metrics.Handler(ctx).(http.HandlerFunc))
	metricsServer := &http.Server{Addr: *metricsListenAddr}

	// 7. Start HTTP server
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
			errc <- httpServer.ListenAndServe()
		}()

		go func() {
			log.Print(ctx, "Metrics server listening", "addr", *metricsListenAddr)
			errc <- metricsServer.ListenAndServe()
		}()

		<-ctx.Done()
		log.Print(ctx, "shutting down HTTP servers")

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpServer.Shutdown(ctx)
		metricsServer.Shutdown(ctx)
	}()

	// Cleanup
	log.Print(ctx, "exiting", "err", <-errc)
	cancel()
	wg.Wait()
	log.Print(ctx, "exited")
}
