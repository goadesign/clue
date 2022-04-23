# log: Smart Logging

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![Go Reference](https://pkg.go.dev/badge/goa.design/clue/log.svg)](https://pkg.go.dev/goa.design/clue/log)

## Overview

Package `log` provides a context-based logging API that intelligently selects
what to log. The API is designed to be used in conjunction with the
[`context`](https://golang.org/pkg/context/) package. The following example
shows how to use the API:

```go
package main

import (
        "context"
        "github.com/goadesign/clue/log"
)

func main() {
        ctx := log.Context(context.Background())
        log.Printf(ctx, "hello %s", "world")
        log.Print(ctx, log.KV{"hello", "world"})

        log.Print(ctx,
                log.KV{"example", "log.KV"},
                log.KV{"order", "deterministic"},
                log.KV{"backed_by", "slice"},
        )

        log.Print(ctx, log.Fields{
                "example": "log.Fields",
                "order": "random",
                "backed_by": "map",
        })
}
```

The example above logs the following messages to stdout (assuming the default
formatter):

```
time=2022-02-22T02:22:02Z level=info msg="hello world"
time=2022-02-22T02:22:02Z level=info hello=world
time=2022-02-22T02:22:02Z level=info example=log.KV order=deterministic backed_by=slice
time=2022-02-22T02:22:02Z level=info order=random backed_by=map example=log.Fields
```

A typical instantiation of the logger for a Goa service looks like this:

```go
ctx := log.With(log.Context(context.Background()), log.KV{"svc", svcgen.ServiceName})
```

Where `svcgen` is the generated Goa service package. This guarantees that all
log entries for the service will have the `svc` key set to the service name.

## Buffering

One of the key features of the `log` package is that it can buffer log messages
until an error occurs (and is logged) or the buffer is explicitely flushed. This
allows the application to write informative log messages without having to worry
about the volume of logs being written.

Specifically any call to the package `Info` function buffers log messages until
the function `Fatal`, `Error` or `Flush` is called. Calls to `Print` are not
buffered. This makes it possible to log specific messages (e.g. request started,
request finished, etc.) without having to flush all the buffered messages.

The following example shows how to use the buffering feature:

```go
log.Infof(ctx, "request started")
// ... no log written so far
log.Errorf(ctx, err, "request failed") // flushes all previous log entries
```

The example above logs the following messages to stdout *after* the call to
`Errorf`:

```
time=2022-02-22T02:22:02Z level=info msg="request started"
time=2022-02-22T02:22:04Z level=error msg="request failed"
```

The `time` key makes it possible to infer the order of log events in case
buffered and non-buffered function calls are mixed. In practice this rarely
happens as non buffered log events are typically created by middlewares which
log before and after the business logic.

### Conditional Buffering

`log` can also be configured to disable buffering conditionally based on the
current context. This makes it possible to force logging when tracing is enabled
for the current request for example.

The following example shows how to conditionally disable the buffering feature:

```go
ctx := log.Context(req.Context(), log.WithDisableBuffering(log.IsTracing))
log.Infof(ctx, "request started") // buffering disabled if tracing is enabled
```

The function given to `WithDisableBuffering` is called with the current context
and should return a boolean indicating whether buffering should be disabled. It
is evaluated upon each call to `Context` and `With`.

### Usage Pattern

Buffering works best in code implementing network request handling (e.g. HTTP or
gRPC requests). The context for each request is used to initialize a new logger
context for example by using the HTTP middleware or gRPC intereceptors defined in
this package (see [below](#http-middleware)). This allows for:

* Creating request specific buffers thereby naturally limiting how many logs are
  kept in memory at a given point in time.
* Evaluating the buffering conditionally based on the request specific context
  (e.g. to disable buffering for traced requests).
* Flushing the buffer when the request encounters an error thereby providing
  useful information about the request.

## Structured Logging

The logging function `Print`, `Debug`, `Info`, `Error` and `Fatal` each accept a
context and a variadic number of key/value pairs. `log` also makes it possible
to build up the log context with a series of key-value pairs via the `With`
function. The following example shows how to leverage structured logging:

```go
ctx := log.Context(context.Background())
ctx := log.With(ctx, log.KV{"key2", "val2"})
log.Print(ctx, log.KV{"hello",  "world 1"})

ctx = log.With(ctx, log.KV{"key3", "val3"})
log.Print(ctx, log.KV{"hello", "world 2"}, log.KV{"key4", "val4"})
```

The example above logs the following message to stdout (assuming the terminal
formatter is being used):

```
INFO[0000] key2=val2 hello="world 1"
INFO[0000] key2=val2 key3=val3 hello="world 2" key4=val4
```

Values must be strings, numbers, booleans, nil or a slice of these types.

## Log Severity

`log` supports three log severities: `debug`, `info`, and `error`. By default
debug logs are not written to the log output. The following example shows how to
enable debug logging:

```go
ctx := log.Context(context.Background())
log.Debugf(ctx, "debug message 1")

ctx := log.Context(ctx, log.WithDebug())
log.Debugf(ctx, "debug message 2")
log.Infof(ctx, "info message")
```

The example above logs the following messages to stdout:

```
DEBG[0000] msg="debug message 2"
INFO[0000] msg="info message"
```

Note that enabling debug logging also disables buffering and causes all future
log messages to be written to the log output as demonstrated above.

## Log Output

By default `log` writes log messages to `os.Stdout`. The following example shows
how to change the log output:

```go
ctx := log.Context(context.Background(), log.WithOutput(os.Stderr))
log.Printf(ctx, "hello world")
```

The example above logs the following message to stderr:

```
INFO[0000] msg="hello world"
```

The `WithOuput` function accepts any type that implements the `io.Writer`
interface.

## Log Format

`log` comes with three predefined log formats and makes it easy to provide
custom formatters. The three built-in formats are:
* `FormatText`: a plain text format using [logfmt](https://brandur.org/logfmt)
* `FormatTerminal`: a format suitable to print logs to colored terminals
* `FormatJSON`: a JSON format

### Text Format

The text format is the default format used when the application is not running
in a terminal.

```go
ctx := log.Context(context.Background(), log.WithFormat(log.FormatText))
log.Printf(ctx, "hello world")
```

The example above logs the following message:

```
time=2022-01-09T20:29:45Z level=info msg="hello world"
```

Where `2022-01-09T20:29:45Z` is the current time in UTC.

### Terminal Format

The terminal format is the default format used when the application is running
in a terminal.

```go
ctx := log.Context(context.Background(), log.WithFormat(log.FormatTerminal))
log.Printf(ctx, "hello world")
```

The example above logs the following message:

```
INFO[0000] msg="hello world"
```

Where `0000` is the number of seconds since the application started. The
severity and each key are colored based on the severity (gray for debug entries,
blue for info entries and red for errors).

### JSON Format

The JSON format prints entries in JSON.

```go
ctx := log.Context(context.Background(), log.WithFormat(log.FormatJSON))
log.Printf(ctx, "hello world")
```

The example above logs the following message:

```
{"time":"2022-01-09T20:29:45Z","level":"info","msg":"hello world"}
```

### Custom Formats

The format can be changed by using the `WithFormat` function as shown above.
Any function that accepts a `Entry` object and returns a slice of bytes can be
used as a format function. The following example shows how to use a custom
format function:

```go
func formatFunc(entry *log.Entry) []byte {
        return []byte(fmt.Sprintf("%s: %s", entry.Severity, entry.Keyvals[0].V))
}

ctx := log.Context(context.Background(), log.WithFormat(formatFunc))
log.Printf(ctx, "hello world")
```

The example above logs the following message to stdout:

```
INFO: hello world
```

## HTTP Middleware

The `log` package includes a HTTP middleware that initializes the request
context with the logger configured in the given context:

```go
check := log.HTTP(ctx)(health.Handler(health.NewChecker(dep1, dep2, ...)))
```

## gRPC Interceptors

The `log` package also includes both unary and stream gRPC interceptor that
initializes the request or stream context with the logger configured in the
given context:

```go
grpcsvr := grpc.NewServer(
	grpcmiddleware.WithUnaryServerChain(
		goagrpcmiddleware.UnaryRequestID(),
		log.UnaryServerInterceptor(ctx),
	))
```

## Standard Logger Compatibility

The `log` package also provides a compatibility layer for the standard
`log` package. The following example shows how to use the standard logger:

```go
ctx := log.Context(context.Background())
logger := log.AsStdLogger(ctx)
logger.Print("hello world")
```

The compatibility layer supports the following functions:
* `Print`
* `Printf`
* `Println`
* `Fatal`
* `Fatalf`
* `Fatalln`
* `Panic`
* `Panicf`
* `Panicln`

The standard logger adapter uses `log.Print` under the hood which means that
there is no buffering when using these functions.

## Goa Request Logging

Loggers created via the `log` package can be adapted to the Goa
[middleware.Logger](https://pkg.go.dev/goa.design/goa/v3/middleware#Logger)
interface. This makes it possible to use this package to configure the logger
used by the middlewares (e.g. to print a log message upon receiving a request
and sending a response).

```go
ctx := log.Context(context.Background())
logger := log.AsGoaMiddlewareLogger(ctx) // logger implements middleware.Logger
```

See the [AsGoaMiddlewareLogger](adapt.go) function for more details on usage.
