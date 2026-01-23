<div align="center">

# clue: Microservice Instrumentation

[![Go Reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/goa.design/clue)
[![License](https://img.shields.io/badge/License-MIT%202.0-blue?style=for-the-badge)](LICENSE)

[![Build Status](https://img.shields.io/github/actions/workflow/status/goadesign/clue/ci.yaml?style=for-the-badge)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![codecov](https://img.shields.io/codecov/c/github/goadesign/clue/main?style=for-the-badge&token=HVP4WT1PS6)](https://codecov.io/gh/goadesign/clue)
[![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=for-the-badge)](https://goreportcard.com/report/goa.design/clue)

</div>

## Overview

Clue provides a set of Go packages for instrumenting microservices. The
emphasis is on simplicity and ease of use. Although not a requirement, Clue
works best when used in microservices written using
[Goa](https://github.com/goadesign/goa).

Clue covers the following topics:

* Instrumentation: the [clue](clue/) package configures OpenTelemetry
  for service instrumentation.
* Logging: the [log](log/) package provides a context-based logging API that
  intelligently selects what to log.
* Health checks: the [health](health/) package provides a simple way for
  services to expose a health check endpoint.
* Dependency mocks: the [mock](mock/) package provides a way to mock
  downstream dependencies for testing.
* Debugging: the [debug](debug/) package makes it possible to troubleshoot
  and profile services at runtime.
* Interceptors: the [interceptors](interceptors/) package provides a set of
  helpful Goa interceptors.

Clue's goal is to provide all the ancillary functionality required to efficiently
operate a microservice style architecture including instrumentation, logging,
debugging and health checks. Clue is not a framework and does not provide
functionality that is already available in the standard library or in other
packages. For example, Clue does not provide a HTTP router or a HTTP server
implementation. Instead, Clue provides a way to instrument existing HTTP or gRPC
servers and clients using the standard library and the OpenTelemetry API.

Learn more about Clue's observability features in the [Observability](https://goa.design/docs/5-real-world/2-observability/)
section of the [goa.design](https://goa.design) documentation. The guide covers how to effectively monitor,
debug and operate microservices using Clue's instrumentation capabilities.

## Packages

* The `clue` package provides a simple API for configuring OpenTelemetry
  instrumentation. The package also provides a way to configure the log
  package to automatically annotate log messages with trace and span IDs.
  The package also implements a dynamic trace sampler that can be used to
  sample traces based on a target maximum number of traces per second.

* The `log` package offers a streamlined, context-based structured logger that
  efficiently buffers log messages. It smartly determines the optimal time to
  flush these messages to the underlying logger. In its default configuration,
  the log package flushes logs upon the logging of an error or when a request is
  traced. This design choice minimizes logging overhead for untraced requests,
  ensuring efficient logging operations.

* The `health` package provides a simple way to expose a health check endpoint
  that can be used by orchestration systems such as Kubernetes to determine
  whether a service is ready to receive traffic. The package implements the
  concept of checkers that can be used to implement custom health checks with
  a default implementation that relies on the ability to ping downstream
  dependencies.

* The `mock` package provides a way to mock downstream dependencies for testing.
  The package provides a tool that generates mock implementations of interfaces
  and a way to configure the generated mocks to validate incoming payloads and
  return canned responses.

* The `debug` package provides a way to dynamically control the log level of a
  running service. The package also provides a way to expose the pprof Go
  profiling endpoints and a way to expose the log level control endpoint.

## Installation

Clue requires Go 1.24 or later. Install the packages required for your
application using `go get`, for example:

```bash
go get goa.design/clue/clue
go get goa.design/clue/log
go get goa.design/clue/health
go get goa.design/clue/mock
go get goa.design/clue/debug
```

## Usage

The following snippet illustrates how to use `clue` to instrument a HTTP server:

```go
ctx := log.Context(context.Background(),    // Create a clue logger context.
    log.WithFunc(log.Span))                 // Log trace and span IDs.

metricExporter := stdoutmetric.New()        // Create metric and span exporters.
spanExporter := stdouttrace.New()
cfg := clue.NewConfig(ctx, "service", "1.0.0", metricExporter, spanExporter)
clue.ConfigureOpenTelemetry(ctx, cfg)       // Configure OpenTelemetry.

handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
})                                          // Create HTTP handler.
handler = otelhttp.NewHandler(handler, "service" ) // Instrument handler.
handler = debug.HTTP()(handler)             // Setup debug log level control.
handler = log.HTTP(ctx)(handler)            // Add logger to request context and log requests.

mux := http.NewServeMux()                   // Create HTTP mux.
mux.HandleFunc("/", handler)                // Mount handler.
debug.MountDebugLogEnabler(mux)             // Mount debug log level control handler.
debug.MountDebugPprof(mux)                  // Mount pprof handlers.
mux.HandleFunc("/health",
    health.NewHandler(health.NewChecker())) // Mount health check handler.

http.ListenAndServe(":8080", mux)           // Start HTTP server.
```

Similarly, the following snippet illustrates how to instrument a gRPC server:

```go
ctx := log.Context(context.Background(),    // Create a clue logger context.
    log.WithFunc(log.Span))                 // Log trace and span IDs.

metricExporter := stdoutmetric.New()
spanExporter := stdouttrace.New()
cfg := clue.NewConfig(ctx, "service", "1.0.0", metricExporter, spanExporter)
clue.ConfigureOpenTelemetry(ctx, cfg)       // Configure OpenTelemetry.

svr := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        log.UnaryServerInterceptor(ctx),    // Add logger to request context and log requests.
        debug.UnaryServerInterceptor()),    // Enable debug log level control
    grpc.StatsHandler(otelgrpc.NewServerHandler()), // Instrument server.
)
```

Note that in the case of gRPC, a separate HTTP server is required to expose the
debug log level control, pprof and health check endpoints:

```go
mux := http.NewServeMux()                   // Create HTTP mux.
debug.MountDebugLogEnabler(mux)             // Mount debug log level control handler.
debug.MountDebugPprof(mux)                  // Mount pprof handlers.
mux.HandleFunc("/health",
    health.NewHandler(health.NewChecker())) // Mount health check handler.

go http.ListenAndServe(":8081", mux)        // Start HTTP server.
```

## Exporters

Exporter are responsible for exporting telemetry data to a backend. The
[OpenTelemetry Exporters documentation](https://opentelemetry.io/docs/instrumentation/go/exporters/)
provides a list of available exporters.

For example, configuring an [OTLP](https://opentelemetry.io/docs/specs/otlp/)
compliant span exporter can be done as follows:

```go
import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
// ...
spanExporter, err := otlptracegrpc.New(
    context.Background(),
    otlptracegrpc.WithEndpoint("localhost:4317"),
    otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
```

While configuring an OTLP compliant metric exporters can be done as follows:

```go
import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
// ...
metricExporter, err := otlpmetricgrpc.New(
    context.Background(),
    otlpmetricgrpc.WithEndpoint("localhost:4317"),
    otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
```

These exporters can then be used to configure Clue:

```go
// Configure OpenTelemetry.
cfg := clue.NewConfig(ctx, "service", "1.0.0", metricExporter, spanExporter)
clue.ConfigureOpenTelemetry(ctx, cfg)
```

## Clients

HTTP clients can be instrumented using the Clue `log` and OpenTelemetry `otelhttptrace` packages. The `log.Client` function wraps a HTTP transport and logs the request and response. The `otelhttptrace.Client` function wraps a HTTP transport and adds OpenTelemetry tracing to the request.

```go
import (
    "context"
    "net/http"
    "net/http/httptrace"

    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttptrace"
    "goa.design/clue/log"
)

// ...
httpc := &http.Client{
    Transport: log.Client(
        otelhttp.NewTransport(
            http.DefaultTransport,
            otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
                return otelhttptrace.NewClientTrace(ctx)
            }),
        ),
    ),
}
```

Similarly, gRPC clients can be instrumented using the Clue `log` and OpenTelemetry `otelgrpc` packages.

```go
import (
    "context"

    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
    "goa.design/clue/log"
)

// ...
grpcconn, err := grpc.DialContext(ctx,
    "localhost:8080",
    grpc.WithUnaryInterceptor(log.UnaryClientInterceptor()),
    grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
)
```

## Custom Instrumentation

Clue relies on the OpenTelemetry API for creating custom instrumentation. The
following snippet illustrates how to create a custom counter and span:

```go
// ... configure OpenTelemetry like in example above
clue.ConfigureOpenTelemetry(ctx, cfg)

// Create a meter and tracer
meter := otel.Meter("mymeter")
tracer := otel.Tracer("mytracer")

// Create a counter
counter, err := meter.Int64Counter("my.counter",
    metric.WithDescription("The number of times the service has been called"),
    metric.WithUnit("{call}"))
if err != nil {
    log.Fatalf("failed to create counter: %s", err)
}

handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Create a span
    ctx, span := tracer.Start(r.Context(), "myhandler")
    defer span.End()

    // ... do something

    // Add custom attributes to span and counter
    attr := attribute.Int("myattr", 42)
    span.SetAttributes(attr)
    counter.Add(ctx, 1, metric.WithAttributes(attr))

    // ... do something else
    if _, err := w.Write([]byte("Hello, World!")); err != nil {
        log.Errorf(ctx, err, "failed to write response")
    }
})
```

## Goa

The `log` package provides a Goa endpoint middleware that adds the service and
method names to the logger context. The `debug` package provides a Goa endpoint
middleware that logs the request and response payloads when debug logging is
enabled. The example below is a snippet extracted from the `main` function of
the
[genforecaster](example/weather/services/forecaster/cmd/forecaster/main.go#L88)
service:

```go
svc := forecaster.New(wc)
endpoints := genforecaster.NewEndpoints(svc)
endpoints.Use(debug.LogPayloads())
endpoints.Use(log.Endpoint)
```

## Example

The [weather](example/weather) example illustrates how to use Clue to instrument
a system of Goa microservices. The example comes with a set of scripts that can
be used to install all necessary dependencies including the
[SigNoz](https://signoz.io/) instrumentation backend.  See the
[README](example/weather/README.md) for more information.

## Migrating from v0.x to v1.x

The v1.0.0 release of `clue` is a major release that introduces breaking
changes. The following sections describe the changes and how to migrate.

### Initialization

The `clue` package provides a cohesive API for both metrics and tracing,
effectively replacing the previous `metrics` and `trace` packages. The
traditional `Context` function, utilized in the `metrics` and `trace` packages
for setting up telemetry, has been deprecated. In its place, the `clue` package
introduces the `NewConfig` function, which generates a `Config` object used in
conjunction with the `ConfigureOpenTelemetry` function to facilitate telemetry
setup.

v0.x:

```go
ctx := log.Context(context.Background())
traceExporter := tracestdout.New()
metricsExporter := metricstdout.New()
ctx = trace.Context(ctx, "service", traceExporter)
ctx = metrics.Context(ctx, "service", metricsExporter)
```

v1.x:

```go
ctx := log.Context(context.Background())
traceExporter := tracestdout.New()
metricsExporter := metricstdout.New()
cfg := clue.NewConfig(ctx, "service", "1.0.0", metricsExporter, traceExporter)
clue.ConfigureOpenTelemetry(ctx, cfg)
```

### Instrumentation

Instrumenting a HTTP service is now done using the standard `otelhttp` package:

v0.x:

```go
handler = trace.HTTP(ctx)(handler)
handler = metrics.HTTP(ctx)(handler)
http.ListenAndServe(":8080", handler)
```

v1.x:

```go
http.ListenAndServe(":8080", otelhttp.NewHandler("service", handler))
```

Similarly, instrumenting a gRPC service is now done using the standard
`otelgrpc` package. v1.x also switches to using a gRPC stats handler instead of
interceptors:

v0.x:

```go
grpcsvr := grpc.NewServer(
        grpcmiddleware.WithUnaryServerChain(
                trace.UnaryServerInterceptor(ctx),
                metrics.UnaryServerInterceptor(ctx),
        ),
        grpcmiddleware.WithStreamServerChain(
                trace.StreamServerInterceptor(ctx),
                metrics.StreamServerInterceptor(ctx),
        ),
)
```

v1.x:

```go
grpcsvr := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
```

### Goa

The `otel` Goa plugin leverages the OpenTelemetry API to annotate spans and
metrics with HTTP routes.

v0.x:

```go
mux := goahttp.NewMuxer()
ctx = metrics.Context(ctx, genfront.ServiceName,
        metrics.WithRouteResolver(func(r *http.Request) string {
                return mux.ResolvePattern(r)
        }),
)
```

v1.x:

```go
package design

import (
        . "goa.design/goa/v3/dsl"
        _ "goa.design/plugins/v3/clue"
)
```

## Contributing

See [Contributing](CONTRIBUTING.md)
