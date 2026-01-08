package log

import (
	"context"
	"io"
	"os"
	"runtime"
	"strconv"

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

	// Output configures where log entries are written and how they are formatted.
	//
	// Output is the unit of configuration for "fanout" logging: a single log entry
	// can be written to multiple outputs, each with its own formatting.
	//
	// Writer and Format must be non-nil.
	Output struct {
		// Writer receives the formatted log bytes.
		Writer io.Writer
		// Format turns a log entry into bytes suitable for Writer.
		Format FormatFunc
	}

	options struct {
		disableBuffering DisableBufferingFunc
		debug            bool
		outputs          []Output
		keyvals          kvList
		kvfuncs          []func(context.Context) []KV
		maxsize          int
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
//
// This option exists for backward compatibility. When used in conjunction with
// WithOutputs, it updates the first output writer.
func WithOutput(w io.Writer) LogOption {
	return func(o *options) {
		if len(o.outputs) == 0 {
			o.outputs = []Output{{Writer: w, Format: FormatText}}
			return
		}
		o.outputs[0].Writer = w
	}
}

// WithFormat sets the log format.
//
// This option exists for backward compatibility. When used in conjunction with
// WithOutputs, it updates the first output format.
func WithFormat(fn FormatFunc) LogOption {
	return func(o *options) {
		if len(o.outputs) == 0 {
			o.outputs = []Output{{Writer: os.Stdout, Format: fn}}
			return
		}
		o.outputs[0].Format = fn
	}
}

// WithOutputs sets the log outputs.
//
// Each output formats the entry then writes it to its writer. This makes it
// possible to log to multiple destinations with independent formats (e.g.
// terminal colors on stdout and JSON to a file).
func WithOutputs(outputs ...Output) LogOption {
	return func(o *options) {
		if len(outputs) == 0 {
			panic("log.WithOutputs: at least one output must be provided")
		}
		for i, out := range outputs {
			if out.Writer == nil {
				panic("log.WithOutputs: output writer is nil")
			}
			if out.Format == nil {
				panic("log.WithOutputs: output format is nil")
			}
			outputs[i] = out
		}
		o.outputs = outputs
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
		outputs:          []Output{{Writer: os.Stdout, Format: format}},
		maxsize:          DefaultMaxSize,
	}
}
