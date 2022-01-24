# log: Smart Logging

[![Build Status](https://github.com/crossnokaye/micro/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/crossnokaye/micro/actions?query=branch%3Amain+event%3Apush)
[![codecov](https://codecov.io/gh/crossnokaye/micro/branch/main/graph/badge.svg?token=HVP4WT1PS6)](https://codecov.io/gh/crossnokaye/micro)

## Overview

Package `log` provides a context-based logging API that intelligently selects
what to log. The API is designed to be used in conjunction with the
[`context`](https://golang.org/pkg/context/) package. The following example
shows how to use the API:

```go
package main

import (
        "context"
        "github.com/crossnokaye/micro/log"
)       

func main() {
        ctx := log.Context(context.Background(), "svc", svcgen.ServiceName)
        ctx := log.With(ctx, "foo", "bar")
        log.Print(ctx, "hello world", "baz", "qux")
}
```

The example above logs the following message to stdout, it is colored if the
application runs in a terminal:

```
INFO[0000] hello world foo=bar baz=qux
```

A typical instantiation of the logger for a Goa service looks like this:

```go
ctx := log.With(log.Context(context.Background()), "svc", svcgen.ServiceName)
```

Where `svcgen` is the generated Goa service package.

## Buffering

One of the key features of the `log` package is that it can buffer log messages
until an error occurs (and is logged) or the buffer is explicitely flushed. This
allows the application to write informative log messages without having to worry
about the volume of logs being written.

Specifically any call to the package `Info` function buffers log messages until
the function `Error` or `Flush` is called. Note that calls to `Print` are not
buffered. This makes it possible to log strategic messages (e.g. request
started, request finished, etc.) without having to flush all the buffered
messages.

The following example shows how to use the buffering feature:

```go
log.Info(ctx, "request started")
// ... no log written so far
log.Error("request failed") // flushes all previous log entries
```

The example above logs the following messages to stdout:

```
INFO[0000] request started
ERRO[0000] request failed
```

### Conditional Buffering

`log` can also be configured to disable buffering conditionally based on the
current context. This makes it possible to force logging when tracing is enabled
for the current request for example.

The following example shows how to conditionally disable the buffering feature:

```go
ctx := log.Context(context.Background(), log.WithDisableBuffering(log.IsTracing))
log.Info(ctx, "request started") // buffering disabled if tracing is enabled
```

The function given to `WithDisableBuffering` is called with the current context
and should return a boolean indicating whether buffering should be disabled. It
is called upon initial creation of the log context and upon each call to `With`.

## Structured Logging

The logging function `Print`, `Debug`, `Info` and `Error` each accept a context,
a message and a variadic number of arguments. The variadic arguments consist of
alternating keys and values. `log` also makes it possible to build up the log
context with a series of key-value pairs via the `With` function and to set
default key-value pairs that should always be logged via the `WithKeyVal`
option. The following example shows how to leverage structured logging:

```go
ctx := log.Context(context.Background(), log.WithKeyValue("key1", "val1"))
ctx := log.With(ctx, "key2", "val2")
log.Print(ctx, "hello world 1")

ctx = log.With(ctx, "key3", "val3")
log.Print(ctx, "hello world 2", "key4", "val4", "key5", "val5")
```

The example above logs the following message to stdout:

```
INFO[0000] hello world 1 key1=val1 key2=val2
INFO[0000] hello world 2 key1=val1 key2=val2 key3=val3 key4=val4 key5=val5
```

Keys of key-value pairs must be strings and values must be strings, numbers,
booleans, nil or a slice of these types.

Log messages are optional and can be omitted by passing an empty string. The
following example shows how to omit a log message:

```go
log.Print(ctx, "", "foo", "bar")
```

The example above logs the following message to stdout:

```
INFO[0000] foo=bar
```

## Log Severity

`log` supports three log severities: `Debug`, `Info`, and `Error`. By default
debug logs are not written to the log output. The following example shows how to
enable debug logging:

```go
ctx := log.Context(context.Background())
log.Debug(ctx, "debug message 1")

ctx := log.Context(ctx, log.WithDebug())
log.Debug(ctx, "debug message 2")
log.Info(ctx, "info message")
```

The example above logs the following messages to stdout:

```
DEBG[0000] debug message 2
INFO[0000] info message
```

Note that enabling debug logging also disables buffering and causes all future
log messages to be written to the log output.

## Log Output

By default `log` writes log messages to `os.Stdout`. The following example shows
how to change the log output:

```go
ctx := log.Context(context.Background(), log.WithOutput(os.Stderr))
log.Print(ctx, "hello world")
```

The example above logs the following message to stderr:

```
INFO[0000] hello world
```

The `WithOuput` function accepts any type that implements the `io.Writer`
interface.

## Log Format

`log` comes with three predefined log formats and makes it easy to provide
custom formatters. The three built-in formats are:
* `FormatText`: a plain text format
* `FormatTerminal`: a format suitable to print logs to colored terminals
* `FormatJSON`: a JSON format

### Text Format

The text format is the default format used when the application is not running
in a terminal.

```go
ctx := log.Context(context.Background(), log.WithFormat(log.FormatText))
log.Print(ctx, "hello world", "foo", "bar")
```

The example above logs the following message:

```
INFO[2022-01-09T20:29:45Z] hello world foo=bar
```

Where `2022-01-09T20:29:45Z` is the current time in UTC.

### Terminal Format

The terminal format is the default format used when the application is running
in a terminal.

```go
ctx := log.Context(context.Background(), log.WithFormat(log.FormatTerminal))
log.Print(ctx, "hello world", "foo", "bar")
```

The example above logs the following message:

```
INFO[0000] hello world foo=bar
```

Where `0000` is the number of seconds since the application started. The
severity and each key are colored based on the severity (gray for debug entries,
blue for info entries and red for errors).

### JSON Format

The JSON format prints entries in JSON.

```go
ctx := log.Context(context.Background(), log.WithFormat(log.FormatJSON))
log.Print(ctx, "hello world", "foo", "bar")
```

The example above logs the following message:

```
{"level":"INFO","time":"2022-01-09T20:29:45Z","message":"hello world","foo":"bar"}
```

### Custom Formats

The format can be changed by using the `WithFormat` function as shown above.
Any function that accepts a `Entry` object and returns a slice of bytes can be
used as a format function. The following example shows how to use a custom
format function:

```go
func formatFunc(entry *log.Entry) []byte {
        return []byte(fmt.Sprintf("%s: %s", entry.Severity, entry.Message))
}

ctx := log.Context(context.Background(), log.WithFormat(formatFunc))
log.Print(ctx, "hello world")
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

## Goa Request Logging

Loggers created via the `log` package can be adapted to the Goa
[middleware.Logger](https://pkg.go.dev/goa.design/goa/v3/middleware#Logger)
interface. This makes it possible to use this package to configure the logger
used by the middlewares (e.g. to print a log message upon receiving a request
and sending a response).

```go
ctx := log.Context(context.Background())
logger := log.Adapt(ctx) // logger implements middleware.Logger
```

See the [Adapt](adapt.go) function for more details on usage.