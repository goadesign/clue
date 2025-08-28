package log

import (
	"context"
	"io"
	"os"
	"runtime"
	"strconv"

	"golang.org/x/term"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
)

type (
	// LogOption is a function that applies a configuration option to a logger.
	LogOption func(*options)

	// DisableBufferingFunc is a function that returns true if the logger
	// should disable buffering for the given context.
	DisableBufferingFunc func(context.Context) bool

	// FormatFunc is a function that formats a log entry.
	FormatFunc func(e *Entry) []byte

	options struct {
		disableBuffering DisableBufferingFunc
		debug            bool
		w                io.Writer
		format           FormatFunc
		keyvals          kvList
		kvfuncs          []func(context.Context) []KV
		maxsize          int
		otellog          otellog.Logger // OpenTelemetry logger
	}
)

// DefaultMaxSize is the default maximum size of a single log message or value
// in bytes. It's also the maximum number of elements in a slice value.
const DefaultMaxSize = 1024

// IsTracing returns true if the context contains a trace created via the
// go.opentelemetry.io/otel/trace package. It is the default
// DisableBufferingFunc used by newly created loggers.
func IsTracing(ctx context.Context) bool {
	span := trace.SpanFromContext(ctx)
	return span.SpanContext().IsValid()
}

// WithDisableBuffering sets the DisableBufferingFunc called to assess whether
// buffering should be disabled.
func WithDisableBuffering(fn DisableBufferingFunc) LogOption {
	return func(o *options) {
		o.disableBuffering = fn
	}
}

// WithDebug enables debug logging and disables buffering.
func WithDebug() LogOption {
	return func(o *options) {
		o.debug = true
	}
}

// WithNoDebug disables debug logging.
func WithNoDebug() LogOption {
	return func(o *options) {
		o.debug = false
	}
}

// WithOutput sets the log output.
func WithOutput(w io.Writer) LogOption {
	return func(o *options) {
		o.w = w
	}
}

// WithFormat sets the log format.
func WithFormat(fn FormatFunc) LogOption {
	return func(o *options) {
		o.format = fn
	}
}

// WithMaxSize sets the maximum size of a single log message or value.
func WithMaxSize(n int) LogOption {
	return func(o *options) {
		o.maxsize = n
	}
}

// WithFileLocation adds the "file" key to each log entry with the parent
// directory, file and line number of the caller: "file=dir/file.go:123".
func WithFileLocation() LogOption {
	return WithFunc(func(context.Context) []KV {
		_, file, line, _ := runtime.Caller(4)
		short := file
		second := false
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				if second {
					short = file[i+1:]
					break
				}
				second = true
			}
		}
		return []KV{{"file", short + ":" + strconv.Itoa(line)}}
	})
}

// WithFunc sets a key/value pair generator function to be called with every
// log entry. The generated key/value pairs are added to the log entry.
func WithFunc(fn func(context.Context) []KV) LogOption {
	return func(o *options) {
		o.kvfuncs = append(o.kvfuncs, fn)
	}
}

// WithOTELLogger sets the OpenTelemetry logger to send logs to OTLP.
func WithOTELLogger(otelLogger otellog.Logger) LogOption {
	return func(o *options) {
		// Safety check: don't set the OTEL logger if it's nil to avoid issues
		if otelLogger != nil {
			o.otellog = otelLogger
		}
	}
}

// IsTerminal returns true if the process is running in a terminal.
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// defaultOptions returns a new options struct with default values.
func defaultOptions() *options {
	format := FormatText
	if IsTerminal() {
		format = FormatTerminal
	}
	return &options{
		disableBuffering: IsTracing,
		w:                os.Stdout,
		format:           format,
		maxsize:          DefaultMaxSize,
	}
}
