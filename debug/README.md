# debug: Inspect Microservices at Runtime

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![Go Reference](https://pkg.go.dev/badge/goa.design/clue/debug.svg)](https://pkg.go.dev/goa.design/clue/debug)

## Overview

Package `debug` provides a few functions to help debug microservices. The
functions make it possible to toggle debug logs on and off at runtime, to log
the content of request and result payloads, and to profile a microservice.

## Usage

### Toggling Debug Logs

The `debug` package provides a `MountDebugLogEnabler` function which adds a
handler to the given mux under the `/debug` path that accepts requests of the
form `/debug?debug-logs=on` and `/debug?debug-logs=off` to manage debug logs
state. The handler returns the current state of debug logs in the response body
and does not change the state if the request does not contain a `debug-logs`
query parameter.  The path, query parameter name and value can be customized by
passing options to the `MountDebugLogEnabler` function.

Note that for the debug log state to take effect, HTTP servers must use handlers
returned by the HTTP function and gRPC servers must make use of the
UnaryInterceptor or StreamInterceptor interceptors.  Also note that gRPC
services must expose an HTTP endpoint to control the debug log state.

```go
// HTTP
mux := http.NewServeMux()
debug.MountDebugLogEnabler(mux)
// ... configure mux with other handlers
srv := &http.Server{Handler: debug.HTTP(mux)}
srv.ListenAndServe()
```

```go
// gRPC
mux := http.NewServeMux()
debug.MountDebugLogEnabler(mux)
srv := &http.Server{Handler: mux}
go srv.ListenAndServe()
gsrv := grpc.NewServer(grpc.UnaryInterceptor(debug.UnaryInterceptor))
lis, _ := net.Listen("tcp", ":8080")
gsrv.Serve(lis)
```

The package also provides a Goa muxer adapter that can be used to mount the
debug log enabler handler on a Goa muxer.

```go
mux := goa.NewMuxer()
debug.MountDebugLogEnabler(debug.Adapt(mux))
```

### Logging Request and Result Payloads

The `debug` package provides a `LogPayloads` Goa endpoint middleware that logs
the content of request and result payloads. The middleware can be used with the
generated `Use` functions on `Endpoints` structs. The middleware is a no-op if
debug logs are disabled.

```go
endpoints := genforecaster.NewEndpoints(svc)
endpoints.Use(debug.LogPayloads())
endpoints.Use(log.Endpoint)
```

### Profiling

The `debug` package provides a `MountPprofHandlers` function which configures a
given mux to serve the pprof handlers under the `/debug/pprof/` URL prefix. The
list of handlers is:

* /debug/pprof/
* /debug/pprof/allocs
* /debug/pprof/block
* /debug/pprof/cmdline
* /debug/pprof/goroutine
* /debug/pprof/heap
* /debug/pprof/mutex
* /debug/pprof/profile
* /debug/pprof/symbol
* /debug/pprof/threadcreate
* /debug/pprof/trace

See the [net/http/pprof](https://pkg.go.dev/net/http/pprof) package
documentation for more information.

The path prefix can be customized by passing an option to the
`MountPprofHandlers` function.

```go
mux := http.NewServeMux()
debug.MountPprofHandlers(mux)
// ... configure mux with other handlers
```

### Example

The weather example illustrates how to make use of this package. In particular
all the services allow for toggling debug logs on or off dynamically (see toggle
HTTP endpoint being mounted
[here](https://github.com/goadesign/clue/blob/main/example/weather/services/forecaster/cmd/forecaster/main.go#L109)
and gRPC interceptor managing log level being
added
[here](https://github.com/goadesign/clue/blob/main/example/weather/services/forecaster/cmd/forecaster/main.go#L89))
and will log the content of request and result payloads when debug
logs are enabled (see Goa endpoint middleware being used
[here](https://github.com/goadesign/clue/blob/main/example/weather/services/forecaster/cmd/forecaster/main.go#L80)).
Additionally the two "back" services (which would not be exposed to the internet
in production) also mount the pprof handlers (see
[here](https://github.com/goadesign/clue/blob/main/example/weather/services/forecaster/cmd/forecaster/main.go#L105)).
With this setup, requests made to the HTTP servers of each service of the form
`/debug?debug-logs=on` turn on debug logs and requests of the form
`/debug?debug-logs=off` turns them back off. Requests made to `/debug/pprof/`
return the pprof package index page while a request to `/debug/pprof/profile`
profile the service for 30s.
