# interceptors: Goa Interceptors

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![Go Reference](https://pkg.go.dev/badge/goa.design/clue/interceptors.svg)](https://pkg.go.dev/goa.design/clue/interceptors)

## Overview

Package `interceptors` provides a set of Goa interceptors that are helpful when
writing microservices.

The following interceptor families are available:

* [Trace Stream](#trace-stream): these interceptor structs and functions allow OpenTelemetry tracing of individual
  messages of a Goa stream whether it is bidirectional, client to server, or server to client. There are
  interceptor structs and functions for use both as client and server interceptors.

## Trace Stream

The Trace Stream family of interceptor structs and functions can be used to trace the individual messages of a Goa stream. This differs from the OpenTelemetry HTTP middleware and gRPC stats handler in that it traces the individual messages of a stream rather than the entire stream which could be long running.

The available Trace Stream interceptor structs are:

* `TraceBidirectionalStreamClientInterceptor`: intercepts a bidirectional stream when used as a client interceptor.
* `TraceServerToClientStreamClientInterceptor`: intercepts a server to client stream when used as a client interceptor.
* `TraceClientToServerStreamClientInterceptor`: intercepts a client to server stream when used as a client interceptor.
* `TraceBidirectionalStreamServerInterceptor`: intercepts a bidirectional stream when used as a server interceptor.
* `TraceServerToClientStreamServerInterceptor`: intercepts a server to client stream when used as a server interceptor.
* `TraceClientToServerStreamServerInterceptor`: intercepts a client to server stream when used as a server interceptor.

The available Trace Stream interceptor functions are:

* `TraceBidirectionalStreamClient`: traces a bidirectional stream when used as a client interceptor.
* `TraceServerToClientStreamClient`: traces a server to client stream when used as a client interceptor.
* `TraceClientToServerStreamClient`: traces a client to server stream when used as a client interceptor.
* `TraceBidirectionalStreamServer`: traces a bidirectional stream when used as a server interceptor.
* `TraceServerToClientStreamServer`: traces a server to client stream when used as a server interceptor.
* `TraceClientToServerStreamServer`: traces a client to server stream when used as a server interceptor.

There are also a set of helper functions that should be used within Goa service method implementations
in order to enable propagation of trace metadata received from streams to a context:

* `SetupTraceStreamRecvContext`: returns a context to be used with the receive method of a stream.
* `GetTraceStreamRecvContext`: returns a context with trace metadata after calling the receive method of a stream.

As a convenience, there are also functions that wrap streams with an interface that handles the work of
calling the `SetupTraceStreamRecvContext` and `GetTraceStreamRecvContext` helper functions:

* `WrapTraceBidirectionalStreamClientStream`: wraps a client stream for a bidirectional stream.
* `WrapTraceServerToClientStreamClientStream`: wraps a client stream for a server to client stream.
* `WrapTraceClientToServerStreamWithResultClientStream`: wraps a client stream for a client to server stream with a result.
* `WrapTraceBidirectionalStreamServerStream`: wraps a server stream for a bidirectional stream.
* `WrapTraceServerToClientStreamServerStream`: wraps a server stream for a server to client stream.
* `WrapTraceClientToServerStreamWithResultServerStream`: wraps a server stream for a client to server stream with a result.

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

var TraceServerToClientStream = Interceptor("TraceServerToClientStream", func() {
    WriteStreamingResult(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
    ReadStreamingResult(func() {
        Attribute("TraceMetadata", MapOf(String, String))
    })
})

var TraceClientToServerStream = Interceptor("TraceClientToServerStream", func() {
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

    Method("MyServerToClientStreamMethod", func() {
        ClientInterceptor(TraceServerToClientStream)
        ServerInterceptor(TraceServerToClientStream)

        ...
    })

    Method("MyClientToServerStreamMethod", func() {
        ClientInterceptor(TraceClientToServerStream)
        ServerInterceptor(TraceClientToServerStream)

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
    return interceptors.TraceBidirectionalStreamClient(ctx, info, next)
}

func (i *MyServiceClientInterceptors) TraceServerToClientStream(ctx context.Context, info *genmyservice.TraceServerToClientStream, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.TraceServerToClientStreamClient(ctx, info, next)
}

func (i *MyServiceClientInterceptors) TraceClientToServerStream(ctx context.Context, info *genmyservice.TraceClientToServerStream, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.TraceClientToServerStreamClient(ctx, info, next)
}

func (i *MyServerServiceInterceptors) TraceBidirectionalStream(ctx context.Context, info *genmyservice.TraceBidirectionalStreamInfo, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.TraceBidirectionalStreamServer(ctx, info, next)
}

func (i *MyServerServiceInterceptors) TraceServerToClientStream(ctx context.Context, info *genmyservice.TraceServerToClientStreamInfo, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.TraceServerToClientStreamServer(ctx, info, next)
}

func (i *MyServerServiceInterceptors) TraceClientToServerStream(ctx context.Context, info *genmyservice.TraceClientToServerStreamInfo, next goa.InterceptorEndpoint) (any, context.Context, error) {
    return interceptors.TraceClientToServerStreamServer(ctx, info, next)
}
```

The interceptor functions take advantage of Go generics to work out of the box with the generated types of your
service as long as you define your interceptors as above.

Alternatively, you can embed the interceptor structs in your interceptor implementations:

```go
import (
    ...
    "goa.design/clue/interceptors"
)

...

type MyServiceClientInterceptors struct {
    interceptors.TraceBidirectionalStreamClientInterceptor[*genmyservice.TraceBidirectionalStreamInfo, genmyservice.MyBidirectionalStreamPayload, genmyservice.MyBidirectionalStreamResult]
    interceptors.TraceServerToClientStreamClientInterceptor[*genmyservice.TraceServerToClientStreamInfo, genmyservice.MyServerToClientStreamResult]
    interceptors.TraceClientToServerStreamClientInterceptor[*genmyservice.TraceClientToServerStreamInfo, genmyservice.MyClientToServerStreamPayload]
}

type MyServerServiceInterceptors struct {
    interceptors.TraceBidirectionalStreamServerInterceptor[*genmyservice.TraceBidirectionalStreamInfo, genmyservice.MyBidirectionalStreamPayload, genmyservice.MyBidirectionalStreamResult]
    interceptors.TraceServerToClientStreamServerInterceptor[*genmyservice.TraceServerToClientStreamInfo, genmyservice.MyServerToClientStreamResult]
    interceptors.TraceClientToServerStreamServerInterceptor[*genmyservice.TraceClientToServerStreamInfo, genmyservice.MyClientToServerStreamPayload]
}
```

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
    ws := interceptors.WrapTraceBidirectionalStreamClientStream(stream)
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
