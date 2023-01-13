# debug: Inspect Microservices at Runtime

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![Go Reference](https://pkg.go.dev/badge/goa.design/clue/debug.svg)](https://pkg.go.dev/goa.design/clue/debug)

## Overview

Package `debug` provides a few functions to help debug microservices. The
functions make it possible to toggle debug logs on and off at runtime, to log
the content of request and result payloads, and to profile a microservice.

## Usage

### Toggling Debug Logs

The `debug` package provides a `MountDebugLogEnabler` function which accepts a
URL prefix, a mux and returns a HTTP middleware. The function configures the mux
to add a handler at the given URL prefix that accepts requests of the form
`/<prefix>?enable=on` and `/<prefix>?enable=off` to turn debug logs on and off.

The function returns a middleware that should be used on all the handlers that
should support toggling debug logs. The package also include unary and stream
gRPC interceptors that should be used on all methods that should support
toggling debug logs.  Note that gRPC services must still expose HTTP endpoints
to enable toggling debug logs, the interceptor merely configures the log level
for each request.

### Logging Request and Result Payloads

The `debug` package provides a `LogPayloads` Goa endpoint middleware that logs
the content of request and result payloads. The middleware can be used with the
generated `Use` functions on `Endpoint` structs. The middleware is a no-op if
debug logs are disabled.

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
`/debug?enable=on` turn on debug logs and requests of the form
`/debug?enable=off` turns them back off. Requests made to `/debug/pprof/` return
the pprof package index page while a request to `/debug/pprof/profile` profile
the service for 30s.