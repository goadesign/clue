# interceptors: Goa Interceptors

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![Go Reference](https://pkg.go.dev/badge/goa.design/clue/interceptors.svg)](https://pkg.go.dev/goa.design/clue/interceptors)

## Overview

Package `interceptors` provides a set of Goa interceptors that are helpful when
writing microservices.

The following interceptor families are available:

* [Trace Stream](#trace-stream): these interceptor functions allow OpenTelemetry tracing of individual
  messages of a Goa stream whether it is bidirectional, client to server, or server to client. There are
  interceptor functions for use both as client and server interceptors.

## Trace Stream

The Trace Stream family of interceptor functions can be used to trace the individual messages of a Goa stream. This differs from the OpenTelemetry HTTP middleware and gRPC stats handler in that it traces the individual messages of a stream rather than the entire stream which could be long running.

The available Trace Stream interceptor functions are:

* `ClientTraceBidirectionalStream`: traces a bidirectional stream when used as a client interceptor.
* `ClientTraceDownStream`: traces a server to client stream when used as a client interceptor.
* `ClientTraceUpStream`: traces a client to server stream when used as a client interceptor.
* `ServerTraceBidirectionalStream`: traces a bidirectional stream when used as a server interceptor.
* `ServerTraceDownStream`: traces a server to client stream when used as a server interceptor.
* `ServerTraceUpStream`: traces a client to server stream when used as a server interceptor.

There are also a set of helper functions that should be used within Goa service method implementations
in order to enable propagation of trace metadata received from streams to a context:

* `SetupTraceStreamRecvContext`: returns a context to be used with the receive method of a stream.
* `GetTraceStreamRecvContext`: returns a context with trace metadata after calling the receive method of a stream.

As a convenience, there are also functions that wrap streams with an interface that handles the work of
calling the `SetupTraceStreamRecvContext` and `GetTraceStreamRecvContext` helper functions:

* `WrapTraceStreamClientBidirectionalStream`: wraps a client stream for a bidirectional stream.
* `WrapTraceStreamClientDownStream`: wraps a client stream for a server to client stream.
* `WrapTraceStreamClientUpStreamWithResult`: wraps a client stream for a client to server stream with a result.
* `WrapTraceStreamServerBidirectionalStream`: wraps a server stream for a bidirectional stream.
* `WrapTraceStreamServerDownStream`: wraps a server stream for a server to client stream.
* `WrapTraceStreamServerUpStreamWithResult`: wraps a server stream for a client to server stream with a result.

These interceptor functions will work best if you also set up OpenTelemetry instrumentation for your service
using the [clue](../clue/) package.

### Usage

In your Goa design, you can define the bidirectional, client to server, and/or server to client Trace Stream
interceptors as follows:

```go
var TraceBidirectionalStream = Interceptor("TraceBidirectionalStream", func() {
    WriteStreamingPayload(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
    ReadStreamingPayload(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
    WriteStreamingResult(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
    ReadStreamingResult(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
})

var TraceDownStream = Interceptor("TraceDownStream", func() {
    WriteStreamingResult(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
    ReadStreamingResult(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
})

var TraceUpStream = Interceptor("TraceUpStream", func() {
    WriteStreamingPayload(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
    ReadStreamingPayload(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
})
```

In the Goa service method definitions where you want to use one of the interceptors, you should specify it as both
a client and server interceptor:

```go
    Method("MyBidirectionalStreamMethod", func() {
        ClientInterceptor(TraceBidirectionalStream)
        ServerInterceptor(TraceBidirectionalStream)

        ...
    })

    Method("MyDownStreamMethod", func() {
        ClientInterceptor(TraceDownStream)
        ServerInterceptor(TraceDownStream)

        ...
    })

    Method("MyUpStreamMethod", func() {
        ClientInterceptor(TraceUpStream)
        ServerInterceptor(TraceUpStream)

        ...
    })
```

In the streaming payload and/or result definitions, you should define the `TraceMetadata` attribute or field:

```go
        Attribute("TraceMetadata", MapOf(String, String))  // for HTTP
        Field(101, "TraceMetadata", MapOf(String, String)) // for gRPC
```

You should generate code from your Goa design as usual using `goa gen`. If you are starting a new service, you can
also use `goa example` to bootstrap it along with examples of the service client and server interceptors.

In your implementation of the service client and server interceptors interfaces, you can call the provided
interceptor functions:

```go
import (
    ...
    "goa.design/clue/interceptors"
)

...

func (i *MyServiceClientInterceptors) TraceBidirectionalStream(ctx context.Context, info *genmyservice.TraceBidirectionalStream, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.ClientTraceBidirectionalStream(ctx, info, next)
}

func (i *MyServiceClientInterceptors) TraceDownStream(ctx context.Context, info *genmyservice.TraceDownStream, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.ClientTraceDownStream(ctx, info, next)
}

func (i *MyServiceClientInterceptors) TraceUpStream(ctx context.Context, info *genmyservice.TraceUpStream, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.ClientTraceUpStream(ctx, info, next)
}

func (i *MyServerServiceInterceptors) TraceBidirectionalStream(ctx context.Context, info *genmyservice.TraceBidirectionalStreamInfo, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.ServerTraceBidirectionalStream(ctx, info, next)
}

func (i *MyServerServiceInterceptors) TraceDownStream(ctx context.Context, info *genmyservice.TraceDownStreamInfo, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.ServerTraceDownStream(ctx, info, next)
}

func (i *MyServerServiceInterceptors) TraceUpStream(ctx context.Context, info *genmyservice.TraceUpStreamInfo, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.ServerTraceUpStream(ctx, info, next)
}
```

The interceptor functions take advantage of Go generics to work out of the box with the generated types of your
service as long as you define your interceptors as above.

Since generated Goa client and server interfaces do not return a context from their receive methods, you will
need to use the helper functions to get the context with the extracted trace metadata after calling the receive
method of the stream:

```go
    ctx = interceptors.SetupTraceStreamRecvContext(ctx, stream)
    result, err := stream.RecvWithContext(ctx)
    ctx = interceptors.GetTraceStreamRecvContext(ctx)
```

Alternatively, you can wrap the stream with an interface that handles the work of calling the
`SetupTraceStreamRecvContext` and `GetTraceStreamRecvContext` helper functions:

```go
    ws := interceptors.WrapTraceStreamClientBidirectionalStream(stream)
    err := ws.Send(ctx, &genmyservice.MyBidirectionalStreamPayload{
        ...
    })
    ...
    ctx, result, err := ws.RecvAndReturnContext(ctx)
    ...
    err = ws.Close()
```

The wrapper functions and interfaces take advantage of Go generics to work out of the box with the generated
types of your service.
