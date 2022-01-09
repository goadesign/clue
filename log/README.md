# log: Smart Logging

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
        ctx := log.Context(context.Background())
        ctx := log.With(ctx, "foo", "bar")
        log.Print(ctx, "hello world")
}
```

The example above logs the following message to stdout:

```
[INFO] [foo=bar] hello world
```

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
ctx := log.Context(context.Background())
log.Info(ctx, "request started")
// ... no log written so far
log.Error("request failed") // flushes all previous log entries
```

The example above logs the following messages to stdout:

```
[INFO] request started
[ERROR] request failed
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
context with a series of key-value pairs via the `With` function. The following
example shows how to leverage structured logging:

```go
ctx := log.Context(context.Background())

ctx := log.With(ctx, "foo", "bar")
log.Print(ctx, "hello world 1")

ctx = log.With(ctx, "baz", "qux")
log.Print(ctx, "hello world 2", "qqux", "qqqux")
```

The example above logs the following message to stdout:

```
[INFO] [foo=bar] hello world 1
[INFO] [foo=bar baz=qux qqux=qqqux] hello world 2
```

Keys of key-value pairs must be strings and values must be strings, numbers,
booleans, nil or a slice of these types.

## Log Levels

`log` supports three log levels: `Debug`, `Info`, and `Error`. By default debug
level logs are not written to the log output. The following example shows how to
enable debug level logging:

```go
ctx := log.Context(context.Background())
log.Debug(ctx, "debug message 1")

ctx := log.Context(ctx, log.WithDebug())
log.Debug(ctx, "debug message 2")
log.Info(ctx, "info message")
```

The example above logs the following messages to stdout:

```
[DEBUG] debug message 2
[INFO] info message
```

Note that setting the log level to `Debug` disables buffering and causes all log
messages to be written to the log output.

## Log Output

By default `log` writes log messages to `os.Stdout`. The following example shows
how to change the log output:

```go
ctx := log.Context(context.Background(), log.WithOutput(os.Stderr))
log.Print(ctx, "hello world")
```

The example above logs the following message to stderr:

```
[INFO] hello world
```

The `WithOuput` function accepts any type that implements the `io.Writer`
interface.

## Log Format

By default `log` writes log messages in the following format:

```
[LEVEL] [key=val key=val ...] message
```

The output is colored if the application is running in a terminal.

The format can be changed by using the `WithFormat` function. The following
example shows how to change the format:

```go
ctx := log.Context(context.Background(), log.WithFormat(log.FormatJSON))
ctx = log.With(ctx, "foo", "bar", "baz", "qux")
log.Print(ctx, "hello world")
```

The example above logs the following message to stdout:

```
{"level":"INFO","message":"hello world","foo":"bar","baz":"qux"}
```

Any function that accepts a `Entry` object and returns a slice of bytes can be
used as a format function. The following example shows how to use a custom
format function:

```go
func formatFunc(entry *log.Entry) []byte {
        return []byte(fmt.Sprintf("%s: %s", entry.Level, entry.Message))
}

ctx := log.Context(context.Background(), log.WithFormat(formatFunc))
log.Print(ctx, "hello world")
```

The example above logs the following message to stdout:

```
INFO: hello world
```