# instrument: Auto Metrics

[![Build Status](https://github.com/crossnokaye/micro/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/crossnokaye/micro/actions?query=branch%3Amain+event%3Apush)
![Coverage](https://img.shields.io/badge/Coverage-93.7%25-brightgreen)

## Overview

Package `instrument` provides a convenient way to add Prometheus metrics to Goa
services. The following example shows how to use the package. It implements an
illustrative `main` function for a fictional service `svc` implemented in the
package `github.com/repo/services/svc`. The service is assumed to expose both
HTTP and gRPC endpoints.

```go
package main

import (
        "context"

        "github.com/crossnokaye/micro/instrument"
        "github.com/crossnokaye/micro/log"
       	goahttp "goa.design/goa/v3/http"

       	"github.com/repo/services/svc"
        httpsvrgen "github.com/repo/services/svc/gen/http/svc/server"
       	grpcsvrgen "github.com/repo/services/svc/gen/grpc/svc/server"
       	svcgen "github.com/repo/services/svc/gen/svc"
)

func main() {
        // Initialize the log context
	ctx := log.With(log.Context(context.Background()), "svc", svcgen.ServiceName)
        // Create the service (user code)
        svc := svc.New(ctx)
        // Wrap the service with Goa endpoints
        endpoints := svcgen.NewEndpoints(svc)

        // Create HTTP server
        mux := goahttp.NewMuxer()
        httpsvr := httpsvrgen.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
        httpsvrgen.Mount(mux, httpsvr)

        // ** Initialize context for instrumentation **
        ctx = instrument.Context(ctx, svcgen.ServiceName)

        // ** Instrument HTTP endpoints **
        handler := instrument.HTTP(ctx)(mux)

        // Create gRPC server
        grpcsvr := grpcsvrgen.New(endpoints, nil)

        // ** Instrument gRPC endpoints **
        unaryInterceptor := instrument.UnaryServerInterceptor(ctx)
        streamInterceptor := instrument.StreamServerInterceptor(ctx)
        pbsvr := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamInterceptor))

        // ** Mount the /metrics handler used by Prometheus to scrape metrics **
        mux.Handle("GET", "/metrics", instrument.Handler())

        // .... Start the servers ....
}
```

The `instrument` functions used to instrument the service are:

* `HTTP`: creates a middleware that instruments an HTTP handler.
* `UnaryServerInterceptor`: creates an interceptor that instruments gRPC unary server methods.
* `StreamServerInterceptor`: creates an interceptor that instruments gRPC stream server methods.
* `Handler`: creates a HTTP handler that exposes Prometheus metrics.

## HTTP Metrics

The middleware returned by the `HTTP` function creates the following metrics:

* `http_server_duration`: Histogram of HTTP request durations in milliseconds.
* `http_server_active_requests`: UpDownCounter of active HTTP requests.
* `http_server_request_size`: Histogram of HTTP request sizes in bytes.
* `http_server_response_size`: Histogram of HTTP response sizes in bytes.

All the metrics have the following labels:

* `goa_service`: The service name as specified in the Goa design.
* `http_verb`: The HTTP verb (`GET`, `POST` etc.).
* `http_host`: The value of the HTTP host header.
* `http_path`: The HTTP path.

All the metrics but `http_server_active_requests` also have the following
additional labels:

* `http_status_code`: The HTTP status code.

## GRPC Metrics

The `UnaryInterceptor` and `StreamInterceptor` functions create the following
metrics:

* `rpc_server_duration`: Histogram of unary request durations in milliseconds.
* `rpc_server_active_requests`: UpDownCounter of active unary and stream requests.
* `rpc_server_request_size`: Histogram of message sizes in bytes, per message for streaming RPCs.
* `rpc_server_response_size`: Histogram of response sizes in bytes, per message for streaming RPCs.
* `rpc_server_stream_message_size`: Histogram of message sizes in bytes, per message for streaming RPCs.
* `rpc_server_stream_response_size`: Histogram of repsonse sizes in bytes, per message for streaming RPCs.

All the metrics have the following labels:

* `goa_service`: The service name as specified in the Goa design.
* `net_peer.addr`: The peer address.
* `rpc_method`: Full name of RPC method.

All the metrics but `rpc_server_active_requests`,
`rpc_server_stream_message_size` and `rpc_rpc_server_stream_response_size` also
have the following additional labels:

* `rpc_status_code`: The response status code.

## Configuration

### Histogram Buckets

The histogram buckets can be specified using the `WithDurationBuckets`,
`WithRequestSizeBuckets` and `WithResponseSizeBuckets` options of the `Context`
function:

```go
ctx = instrument.Context(ctx, svc.ServiceName,
        instrument.WithDurationBuckets([]float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}),
        instrument.WithRequestSizeBuckets([]float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}),
        instrument.WithResponseSizeBuckets([]float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}))
```

The default bucket upper boundaries are:

* `10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, +Inf` for request duration.
* `10, 100, 500, 1000, 5000, 10000, 50000, 100000, 1000000, 10000000, +Inf` for request and response size.

### Prometheus Registry

By default `instrument` uses the global Prometheus registry to create the
metrics and serve them. A user configured registerer can be specified when
creating the metrics via `WithRegisterer`:

```go
ctx = instrument.Context(ctx, svc.ServiceName, instrument.WithRegisterer(registerer))(mux)
```

A user configured gatherer (used to collect the metrics) and registerer (used to
register metrics for the `/metrics` endpoint) can be specified when creating the
metrics handler via `WithGatherer` and `WithHandlerRegisterer` respectively:

```go
handler = instrument.Handler(ctx, instrument.WithGatherer(gatherer), instrument.WithHandlerRegisterer(registerer))
```
