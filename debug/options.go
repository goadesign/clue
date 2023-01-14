package debug

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	// LogPayloadsOption is a function that applies a configuration option
	// to the LogPayloads middleware.
	LogPayloadsOption func(*options)

	// FormatFunc is used to format the logged value for payloads and
	// results.
	FormatFunc func(context.Context, interface{}) string

	options struct {
		maxsize int // maximum number of bytes in a single log message or value
		format  FormatFunc
		client  bool
	}
)

// DefaultMaxSize is the default maximum size for a logged request or result
// value in bytes.
const DefaultMaxSize = 1024

// WithFormat sets the log format.
func WithFormat(fn FormatFunc) LogPayloadsOption {
	return func(o *options) {
		o.format = fn
	}
}

// WithMaxSize sets the maximum size of a single log message or value.
func WithMaxSize(n int) LogPayloadsOption {
	return func(o *options) {
		o.maxsize = n
	}
}

// WithClient prefixes the log keys with "client-". This is useful when
// logging client requests and responses.
func WithClient() LogPayloadsOption {
	return func(o *options) {
		o.client = true
	}
}

// FormatJSON returns a function that formats the given value as JSON.
func FormatJSON(ctx context.Context, v interface{}) string {
	js, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("<invalid: %s>", err)
	}
	return string(js)
}

// defaultOptions returns a new options struct with default values.
func defaultOptions() *options {
	return &options{
		format:  FormatJSON,
		maxsize: DefaultMaxSize,
	}
}
