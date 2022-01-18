# tracing: Simple Request Tracing

[![Build Status](https://github.com/crossnokaye/micro/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/crossnokaye/micro/actions?query=branch%3Amain+event%3Apush)
[![codecov](https://codecov.io/gh/crossnokaye/micro/branch/main/graph/badge.svg?token=HVP4WT1PS6)](https://codecov.io/gh/crossnokaye/micro)

## Overview

Package `tracing` provides request tracing functionality that follows the
[OpenTelemetry](https://opentelemetry.io/) specification. The package is
designed to be used in conjunction with [Goa](https://goa.design/). In
particular it is aware of the Goa RequestID
[HTTP](https://github.com/goadesign/goa/blob/v3/http/middleware/requestid.go)
and
[gRPC](https://github.com/goadesign/goa/blob/v3/grpc/middleware/requestid.go)
middlewares and adds both the Goa service name and current request ID to the
span attributes.

The implementation makes use of the 
[AWS Distro for OpenTelemetry Go SDK](https://github.com/aws-observability/aws-otel-go)
and is therefore compatible with the AWS X-Ray service.

## Usage

The following example shows how to use the package. It implements an
illustrative `main` function for a fictional service `svc` implemented in the
package `github.com/repo/services/svc`

```go
package main
        
import (
       "context"

       "github.com/crossnokaye/micro/log"
       "github.com/crossnokaye/micro/tracing"
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

        // ** Create trace provider **
        collectorAddr := "localhost:6831" // AWS X-Ray collector address
        provider := tracing.NewTraceProvider(ctx, svcgen.ServiceName, collectorAddr)

        // ** Trace HTTP requests **
        handler := tracing.HTTP(ctx, provider)(mux)

        // Create gRPC server
        grpcsvr := grpcsvrgen.New(endpoints, nil)

        // ** Trace gRPC requests **
        u := tracing.UnaryServerTrace(provider)
        s := tracing.StreamServerTrace(provider)
        pbsvr := grpc.NewServer(grpc.UnaryInterceptor(u), grpc.StreamInterceptor(s))

        // ...
}
```

The tracing package uses an adaptative sampler that is configured to sample at a
given maximum number of request per seconds (2 per default). Using a time based
rate rather than e.g. a fixed percentage rate allows the sampler to adapt to the
load of the service.

### Making Requests to Downstream Dependencies

For tracing to work appropriately the tracing package must be used when making
requests to downstream dependencies. 

For HTTP dependencies the tracing package provides a `WrapDoer` function that
can be used to wrap a `http.Client` to trace all requests made through it. The
implementationg firsts validates that the current request is being traced and
only adds a span if a trace is active. Example:

```go
// Create a tracing HTTP client
doer := tracing.WrapDoer(http.DefaultClient)
```

For gRPC dependencies the tracing package provides the `UnaryClientTrace` and
`StreamClientTrace` interceptors that can be used when making gRPC calls. These
functions will create a span for the current request if it is traced. Example:

```go
// Create a tracing client for gRPC unary calls
conn, err := grpc.Dial(url, grpc.WithUnaryInterceptor(UnaryClientTrace(provider)))

// Create a tracing client for gRPC stream calls
conn, err := grpc.Dial(url, grpc.WithStreamInterceptor(StreamClientTrace(provider)))
```

### Creating Additional Spans

Once configured the tracing package will automatically create spans for a sample
of incoming requests. The function `IsTraceActive` can be used to determine if
the current request is being traced.

The tracing package also provides a `Start` function that can be used to create a
new span. `Start` only creates a space for requests with an active trace.

```go
func (s *svc) DoSomething(ctx context.Context, req *svcgen.DoSomethingRequest) (*svcgen.DoSomethingResponse, error) {
        // ...
        // Create a child span to measure the time taken to run an intensive
        // operation.
        span := tracing.Start(ctx, "DoSomethingIntense")
        DoSomethingIntense(ctx)
        span.End()
        // ...
}
```



