# instrument: Auto Metrics

[![Build Status](https://github.com/crossnokaye/micro/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/crossnokaye/micro/actions?query=branch%3Amain+event%3Apush)
![Coverage](https://img.shields.io/badge/Coverage-93.7%25-brightgreen)

## Overview

Package `instrument` provides a convenient way to add Prometheus metrics to Goa
services. The following example shows how to use the package. It implements an
illustrative `main` function for a fictional service `svc` implemented in the
package `github.com/repo/services/svc`. The service is assumed to expose both


```go
package main

import (
        "context"

        "github.com/crossnokaye/micro/instrument"
        "github.com/crossnokaye/micro/log"
       	goahttp "goa.design/goa/v3/http"

       	"github.com/repo/services/svc"
        httpserver "github.com/repo/services/svc/gen/http/svc/server"
       	grpcserver "github.com/repo/services/svc/gen/grpc/svc/server"
       	svcgen "github.com/repo/services/svc/gen/svc"
)

func main() {
        // Initialize the log context
	ctx := log.With(log.Context(context.Background(), log.WithFormat(format)), "svc", svcgen.ServiceName)
        // Create the service (user code)
        svc := svc.New(ctx)
        // Wrap the service with Goa endpoints
        endpoints := svcgen.NewEndpoints(svc)

        // Create HTTP server
        mux := goahttp.NewMuxer()
        httpsvr := httpserver.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
        httpserver.Mount(mux, httpsvr)

        // ** Instrument HTTP handler **
        handler := instrument.HTTP(svcgen.ServiceName)(mux)

        // Create gRPC server
        grpcsvr := grpcserver.New(endpoints, nil)

        // ** Instrument gRPC server **
        unaryInterceptor, err := instrument.UnaryServerInterceptor(ctx, svcgen.ServiceName)
        if err != nil {
                log.Error(ctx, err)
        }
        streamInterceptor, err := instrument.StreamServerInterceptor(ctx, svcgen.ServiceName)
        if err != nil {
                log.Error(ctx, err)
        }
        pbsvr := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamInterceptor))

        // ** Mount the /metrics handler used by Prometheus to scrape metrics **
        mux.Handle("GET", "/metrics", instrument.Handler())

        // .... Start the servers ....
}
```

The `instrument` functions used to instrument the service are:

* `HTTP`: creates a middlware that instruments an HTTP handler.
* `UnaryServerInterceptor`: creates an interceptor that instruments gRPC unary server methods.
* `StreamServerInterceptor`: creates an interceptor that instruments gRPC stream server methods.
* `Handler`: creates a HTTP handler that exposes Prometheus metrics.

## HTTP Metrics

The middleware returned by the `HTTP` function creates the following metrics:

* `http.server.duration`: Histogram of HTTP request durations in milliseconds.
* `http.server.active_requests`: UpDownCounter of active HTTP requests.
* `http.server.request.size`: Histogram of HTTP request sizes in bytes.
* `http.server.response.size`: Histogram of HTTP response sizes in bytes.

All the metrics have the following labels:

* `goa.service`: The service name as specified in the Goa design.
* `http.verb`: The HTTP verb (`GET`, `POST` etc.).
* `http.host`: The value of the HTTP host header.
* `http.path`: The HTTP path.

All the metrics but `http.server.active_requests` also have the following
additional labels:

* `goa.method`: The method name as specified in the Goa design.
* `http.status_code`: The HTTP status code.

## GRPC Metrics

The `UnaryInterceptor` and `StreamInterceptor` functions create the following
metrics:

* `rpc.server.duration`: Histogram of unary request durations in milliseconds.
* `rpc.server.active_requests`: UpDownCounter of active unary and stream requests.
* `rpc.server.request.size`: Histogram of message sizes in bytes, per message for streaming RPCs.
* `rpc.server.response.size`: Histogram of response sizes in bytes, per message for streaming RPCs.

All the metrics have the following labels:

* `goa.service`: The service name as specified in the Goa design.
* `net.peer.addr`: The peer address.
* `rpc.method`: Full name of RPC method.

All the metrics but `rpc.server.active_requests` also have the following
additional labels:

* `goa.method`: The method name as specified in the Goa design.
* `rpc.status_code`: The response status code.

## Configuration

### Histogram Buckets
The histogram buckets can be specified using the `WithDurationBuckets`,
`WithRequestSizeBuckets` and `WithResponseSizeBuckets` options:

```go
err := instrument.HTTP(svc.ServiceName,
        instrument.WithDurationBuckets([]float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}),
        instrument.WithRequestSizeBuckets([]float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}),
        instrument.WithResponseSizeBuckets([]float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}),
        )(mux)
```

The default bucket upper boundaries are:

* `10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, +Inf` for request duration.
* `10, 100, 500, 1000, 5000, 10000, 50000, 100000, 1000000, 10000000, +Inf` for request and response size.

### Prometheus Registry

By default `instrument` uses the global Prometheus registry to create the
metrics and serve them. A user configured registerer can be specified when
creating the metrics via `WithRegisterer`:

```go
err := instrument.HTTP(svc.ServiceName, instrument.WithRegisterer(registerer))(mux)
```

A user configured gatherer (used to collect the metrics) and registerer (used to
register metrics for the `/metrics` endpoint) can be specified when creating the
metrics handler via `WithGatherer` and `WithHandlerRegisterer` respectively:

```go
err := instrument.Handler(ctx, instrument.WithGatherer(gatherer), instrument.WithHandlerRegisterer(registerer))
```
