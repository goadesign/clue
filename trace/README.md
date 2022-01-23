# trace: Simple Request Tracing

[![Build Status](https://github.com/crossnokaye/micro/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/crossnokaye/micro/actions?query=branch%3Amain+event%3Apush)
[![codecov](https://codecov.io/gh/crossnokaye/micro/branch/main/graph/badge.svg?token=HVP4WT1PS6)](https://codecov.io/gh/crossnokaye/micro)

## Overview

Package `trace` provides request tracing functionality that follows the
[OpenTelemetry](https://opentelemetry.io/) specification. The package is
designed to be used in conjunction with [Goa](https://goa.design/). In
particular it is aware of the Goa RequestID
[HTTP](https://github.com/goadesign/goa/blob/v3/http/middleware/requestid.go)
and
[gRPC](https://github.com/goadesign/goa/blob/v3/grpc/middleware/requestid.go)
middlewares and adds both the Goa service name and current request ID to the
span attributes.

The package uses an adaptative sampler that is configured to sample at a given
maximum number of request per seconds (2 per default). Using a time based rate
rather than e.g. a fixed percentage rate allows the sampler to adapt to the load
of the service.

## Usage

The following example shows how to use the package. It implements an
illustrative `main` function for a fictional service `svc` implemented in the
package `github.com/repo/services/svc`

```go
package main
        
import (
       "context"

       "github.com/crossnokaye/micro/log"
       "github.com/crossnokaye/micro/trace"
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

        // ** Initialize context for tracing **
        conn, err := grpc.DialContext(ctx, "localhost:6831") // Remote span collector address
        if err != nil {
                log.Error(ctx, "unable to connect to span collector", "err", err)
                os.Exit(1)
        }
        ctx = trace.Context(ctx, svcgen.ServiceName, conn)

        // ** Trace HTTP requests **
        handler := trace.HTTP(ctx, svcgen.ServiceName)(mux)

        // Create gRPC server
        grpcsvr := grpcsvrgen.New(endpoints, nil)

        // ** Trace gRPC requests **
        u := trace.UnaryServerTrace(ctx)
        s := trace.StreamServerTrace(ctx)
        pbsvr := grpc.NewServer(grpc.UnaryInterceptor(u), grpc.StreamInterceptor(s))

        // ...
}
```

### Making Requests to Downstream Dependencies

For tracing to work appropriately all clients to downstream dependencies must be
configured using the appropriate trace package function. 

For HTTP dependencies the trace package provides a `TraceClient` function that
can be used to configure a `http.Client` to trace all requests made through it.
`TraceClient` does nothing if the current request is not traced or the context
not initialized with `trace.Context`.

```go
// Create a tracing HTTP client
doer := trace.WrapDoer(ctx, http.DefaultClient)
```

For gRPC dependencies the trace package provides the `UnaryClientTrace` and
`StreamClientTrace` interceptors that can be used when making gRPC calls. These
functions will create a span for the current request if it is traced. Example:

```go
// Create a tracing client for gRPC unary calls
conn, err := grpc.Dial(url, grpc.WithUnaryInterceptor(UnaryClientTrace(ctx)))

// Create a tracing client for gRPC stream calls
conn, err := grpc.Dial(url, grpc.WithStreamInterceptor(StreamClientTrace(ctx)))
```

### Creating Additional Spans

Once configured the trace package automatically creates spans for a sample of
incoming requests. The function `IsTraced` can be used to determine if the
current request is being traced.

The trace package also provides a `StartSpan` function that can be used to
create a new child span. The caller must also call `EndSpan` when the span is
complete. Both functions do nothing if the current request is not being traced
or the context has not been initialized with `trace.Context`.

```go
func (s *svc) DoSomething(ctx context.Context, req *svcgen.DoSomethingRequest) (*svcgen.DoSomethingResponse, error) {
        // ...
        // Create a child span to measure the time taken to run an intensive
        // operation.
        ctx = trace.StartSpan(ctx, "DoSomethingIntense")
        DoSomethingIntense(ctx)
        trace.EndSpan(ctx)
        // ...
}
```

### Adding Attributes to Spans

Attributes decorate spans and add contextual information to the trace. By default
this package adds the following attributes to HTTP requests:

* `http.scheme`: The HTTP scheme (`http` or `https`).
* `http.host`: The host name of the request if available in the `Host` field.
* `http.flavor`: The flavor of the request (`1.0`, `1.1` or `2`).
* `http.method`: The HTTP method of the request.
* `http.target`: The URI of the request.
* `http.client_id`: The client IP present in the `X-Forwarded-For` header if any.
* `http.user_agent`: The user agent of the request if any.
* `http.request_content_length`: The length of the request body if any.
* `enduser.id`: The request basic auth username if any.
* `net.transport`: One of `ip_tcp`, `ip_udp`, `ip`, `unix` or `other`.
* `net.peer.ip`, `net.peer.name`, `net.peer.port`: The IP address, port and name 
  of the remote peer if available in the request `RemoteAddr` field.
* `net.host.ip`, `net.host.name`, `net.host.port`: The IP address, port and name
  of the remote host if available in the request `Host` field, the request `Host`
  header or the request URL.

and the following attributes to gRPC requests:

* `rpc.system`: always set to `grpc`.
* `rpc.service`: The gRPC service name.
* `rpc.method`: The gRPC method name.
* `net.peer.ip`, `net.peer.port`: The IP address and port of the remote peer.

Service method logic can add attributes when creating new spans via the
`WithAttributes` option. Custom attributes can also be added later on with
`SetSpanAttributes`.  `SetSpanAttributes` does nothing if the request is not
traced or the context not initialized with `trace.Context`.

```go
// Create a child span with attributes
ctx = trace.StartSpan(ctx, "DoSomething", trace.WithAttributes(
        "key1", "value1",
        "key2", "value2",
))

// Add a custom attribute to the current span
trace.SetSpanAttributes(ctx, "custom_attribute", "value")
```

### Adding Events

The `AddEvent` function makes it possible to attach events to a span. Events are
useful to trace operations that are too fast to have their own span. For
example, the completion of an asynchronous operation. Attributes can be added to
the event to add contextual information. `AddEvent` does nothing if the request
is not traced or the context not initialized with `trace.Context`.

```go
// Add an event to the current span
trace.AddEvent(ctx, "operation completed", "operation_id", operationID, "status", status) 
```

### Span Status And Error

The `Succeed`, `Fail` and `RecordError` functions can be used to set the status
and error of a span. The status indicates the success or failure of the
operation.  The error is used to record the error that occurred if any.  Note
that recording an error does not automatically change the status of the span.
The functions do nothing if the request is not traced or the context not
initialized with `trace.Context`. The default status of a completed span is
success.

```go
// Set the status of the current span to success (default)
trace.Succeed(ctx)

// Record an error in the current span and set the status to failure
trace.RecordError(ctx, err)
trace.Fail(ctx, "operation failed")
```

