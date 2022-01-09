package log

import (
	"context"
	"io"
	"os"
	"syscall"

	"golang.org/x/term"

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
	}
)

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

// defaultOptions returns a new options struct with default options.
func defaultOptions() *options {
	format := FormatText
	if term.IsTerminal(syscall.Stdin) {
		format = FormatTerminal
	}
	return &options{
		disableBuffering: IsTracing,
		w:                os.Stdout,
		format:           format,
	}
}
