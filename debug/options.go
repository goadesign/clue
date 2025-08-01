package debug

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	// LogPayloadsOption is a function that applies a configuration option
	// to the LogPayloads middleware.
	LogPayloadsOption func(*lpOptions)

	// DebugLogEnablerOption is a function that applies a configuration option
	// to MountDebugLogEnabler.
	DebugLogEnablerOption func(*dleOptions)

	// PprofOption is a function that applies a configuration option to
	// MountPprofHandlers.
	PprofOption func(*pprofOptions)

	// FormatFunc is used to format the logged value for payloads and
	// results.
	FormatFunc func(context.Context, any) string

	lpOptions struct {
		maxsize int // maximum number of bytes in a single log message or value
		format  FormatFunc
		client  bool
	}

	dleOptions struct {
		path   string
		query  string
		onval  string
		offval string
	}

	pprofOptions struct {
		prefix string
	}
)

// DefaultMaxSize is the default maximum size for a logged request or result
// value in bytes used by LogPayloads.
const DefaultMaxSize = 1024

// WithFormat sets the log format used by LogPayloads.
func WithFormat(fn FormatFunc) LogPayloadsOption {
	return func(o *lpOptions) {
		o.format = fn
	}
}

// WithMaxSize sets the maximum size of a single log message or value used by
// LogPayloads.
func WithMaxSize(n int) LogPayloadsOption {
	return func(o *lpOptions) {
		o.maxsize = n
	}
}

// WithClient prefixes the log keys used by LogPayloads with "client-". This is
// useful when logging client requests and responses.
func WithClient() LogPayloadsOption {
	return func(o *lpOptions) {
		o.client = true
	}
}

// WithPath sets the URL path used by MountDebugLogEnabler.
func WithPath(path string) DebugLogEnablerOption {
	return func(o *dleOptions) {
		o.path = path
	}
}

// WIthPrefix sets the path prefix used by MountPprofHandlers.
func WithPrefix(prefix string) PprofOption {
	return func(o *pprofOptions) {
		o.prefix = prefix
	}
}

// WithQuery sets the query string parameter name used by MountDebugLogEnabler
// to enable or disable debug logs.
func WithQuery(query string) DebugLogEnablerOption {
	return func(o *dleOptions) {
		o.query = query
	}
}

// WithOnValue sets the query string parameter value used by
// MountDebugLogEnabler to enable debug logs.
func WithOnValue(onval string) DebugLogEnablerOption {
	return func(o *dleOptions) {
		o.onval = onval
	}
}

// WithOffValue sets the query string parameter value used by
// MountDebugLogEnabler to disable debug logs.
func WithOffValue(offval string) DebugLogEnablerOption {
	return func(o *dleOptions) {
		o.offval = offval
	}
}

// FormatJSON returns a function that formats the given value as JSON.
func FormatJSON(_ context.Context, v any) string {
	js, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("<invalid: %s>", err)
	}
	return string(js)
}

// defaultLogPayloadsOptions returns a new lpOptions struct with default values.
func defaultLogPayloadsOptions() *lpOptions {
	return &lpOptions{
		format:  FormatJSON,
		maxsize: DefaultMaxSize,
	}
}

// defaultDebugLogEnablerOptions returns a new dleOptions struct with default values.
func defaultDebugLogEnablerOptions() *dleOptions {
	return &dleOptions{
		path:   "debug",
		query:  "debug-logs",
		onval:  "on",
		offval: "off",
	}
}

// defaultPprofOptions returns a new pprofOptions struct with default values.
func defaultPprofOptions() *pprofOptions {
	return &pprofOptions{
		prefix: "/debug/pprof/",
	}
}
