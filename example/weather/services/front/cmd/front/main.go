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

	"github.com/crossnokaye/micro/health"
	"github.com/crossnokaye/micro/log"
	"github.com/crossnokaye/micro/trace"
	goahttp "goa.design/goa/v3/http"
	goahttpmiddleware "goa.design/goa/v3/http/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/crossnokaye/micro/example/weather/services/front"
	"github.com/crossnokaye/micro/example/weather/services/front/clients/forecast"
	"github.com/crossnokaye/micro/example/weather/services/front/clients/locator"
	genfront "github.com/crossnokaye/micro/example/weather/services/front/gen/front"
	genhttp "github.com/crossnokaye/micro/example/weather/services/front/gen/http/front/server"
)

func main() {
	var (
		httpListenAddr     = flag.String("http", ":8084", "HTTP listen address (health checks)")
		forecastAddr       = flag.String("forecast", ":8080", "Forecast service address")
		forecastHealthAddr = flag.String("forecast-health", ":8081", "Forecast service health-check address")
		locatorAddr        = flag.String("locator", ":8082", "Locator service address")
		locatorHealthAddr  = flag.String("locator-health", ":8083", "Locator service health-check address")
		collectorAddr      = flag.String("coladdr", ":55681", "OpenTelemetry remote collector address")
		debug              = flag.Bool("debug", false, "Enable debug logs")
	)
	flag.Parse()

	// 1. Create logger
	ctx := log.With(log.Context(context.Background()), "svc", genfront.ServiceName)
	if *debug {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debug(ctx, "debug logs enabled")
	}

	// 2. Setup tracing
	conn, err := grpc.DialContext(ctx, *collectorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(ctx, "failed to connect to OpenTelementry collector", "err", err)
		os.Exit(1)
	}
	ctx, err = trace.Context(ctx, genfront.ServiceName, conn)
	if err != nil {
		log.Error(ctx, "failed to initialize tracing", "err", err)
		os.Exit(1)
	}

	// 3. Create clients
	lcc, err := grpc.DialContext(ctx, *locatorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(ctx, "failed to connect to locator", "err", err)
		os.Exit(1)
	}
	lc := locator.New(lcc)
	fcc, err := grpc.DialContext(ctx, *forecastAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(ctx, "failed to connect to forecast", "err", err)
		os.Exit(1)
	}
	fc := forecast.New(fcc)

	// 4. Create service & endpoints
	svc := front.New(fc, lc)
	endpoints := genfront.NewEndpoints(svc)
	endpoints.Use(log.SetContext(ctx))

	// 5. Create transport
	mux := goahttp.NewMuxer()
	server := genhttp.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	genhttp.Mount(mux, server)
	handler := goahttpmiddleware.Log(log.Adapt(ctx))(mux)
	handler = goahttpmiddleware.RequestID()(handler)
	l := &http.Server{Addr: *httpListenAddr, Handler: handler}
	for _, m := range server.Mounts {
		log.Print(ctx, "mount", "method", m.Method, "verb", m.Verb, "path", m.Pattern)
	}

	// 6. Mount health check
	check := health.Handler(health.NewChecker(
		health.NewPinger("locator", "http", *locatorHealthAddr),
		health.NewPinger("forecast", "http", *forecastHealthAddr)))
	check = log.HTTP(ctx)(check).(http.HandlerFunc)
	mux.Handle("GET", "/healthz", check)
	mux.Handle("GET", "/livez", check)

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
			errc <- l.ListenAndServe()
		}()

		<-ctx.Done()
		log.Print(ctx, "shutting down HTTP server")

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		l.Shutdown(ctx)
	}()

	// Cleanup
	log.Print(ctx, "exiting", "err", <-errc)
	cancel()
	wg.Wait()
	log.Print(ctx, "exited")
}
