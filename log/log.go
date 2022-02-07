package log

import (
	"context"
	"sync"
	"time"
)

type (
	// KeyVals represents a list of key/value pairs.
	KeyVals []interface{}

	// Log entry
	Entry struct {
		Time     time.Time
		Severity Severity
		KeyVals  KeyVals
		Message  string
	}

	// Logger implementation
	logger struct {
		options *options
		lock    sync.Mutex
		keyvals []interface{}
		entries []*Entry
		flushed bool
	}

	// Log severity enum
	Severity int

	// private type for context keys
	ctxKey int
)

const (
	SeverityDebug Severity = iota + 1
	SeverityInfo
	SeverityError
)

const (
	ctxLogger ctxKey = iota + 1
)

// Be kind to tests
var timeNow = time.Now

// Context initializes a context for logging.
func Context(ctx context.Context, opts ...LogOption) context.Context {
	var l *logger
	if v := ctx.Value(ctxLogger); v != nil {
		l = v.(*logger)
	} else {
		l = &logger{options: defaultOptions()}
	}
	for _, opt := range opts {
		opt(l.options)
	}
	return context.WithValue(ctx, ctxLogger, l)
}

// Debug logs a debug message. msg is optional and can be empty. keyvals is an
// alternating list of keys and values. Keys must be strings and values must be
// strings, numbers, booleans, nil or a slice of these types.
func Debug(ctx context.Context, msg string, keyvals ...interface{}) {
	log(ctx, SeverityDebug, true, msg, keyvals...)
}

// Print logs an info message and ignores buffering. msg is optional and can be
// empty. keyvals is an alternating list of keys and values. Keys must be
// strings and values must be strings, numbers, booleans, nil or a slice of
// these types.
func Print(ctx context.Context, msg string, keyvals ...interface{}) {
	log(ctx, SeverityInfo, false, msg, keyvals...)
}

// Info logs an info message. msg is optional and can be empty. keyvals is an
// alternating list of keys and values. Keys must be strings and values must be
// strings, numbers, booleans, nil or a slice of these types.
func Info(ctx context.Context, msg string, keyvals ...interface{}) {
	log(ctx, SeverityInfo, true, msg, keyvals...)
}

// Error logs an error message and flushes the log buffer if not already
// flushed. msg is optional and can be empty. keyvals is an alternating list of
// keys and values. Keys must be strings and values must be strings, numbers,
// booleans, nil or a slice of these types.
func Error(ctx context.Context, msg string, keyvals ...interface{}) {
	Flush(ctx)
	log(ctx, SeverityError, true, msg, keyvals...)
}

// With creates a copy of the given log context and appends the given key/value
// pairs to it. keyvals is an alternating list of keys and values. Keys must be
// strings and values must be strings, numbers, booleans, nil or a slice of
// these types.
func With(ctx context.Context, keyvals ...interface{}) context.Context {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, nil)
	}
	v := ctx.Value(ctxLogger)
	if v == nil {
		return ctx
	}
	l := v.(*logger)
	l.lock.Lock()
	copy := logger{
		options: l.options,
		keyvals: l.keyvals,
		entries: l.entries,
		flushed: l.flushed,
	}
	l.lock.Unlock()
	copy.keyvals = append(copy.keyvals, keyvals...)

	// Make sure that if Go needs to grow the slice then each context get
	// its own memory.
	copy.keyvals = copy.keyvals[:len(copy.keyvals):len(copy.keyvals)]
	copy.entries = copy.entries[:len(copy.entries):len(copy.entries)]

	return context.WithValue(ctx, ctxLogger, &copy)
}

// Flush flushes the log entries to the writer.
func Flush(ctx context.Context) {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return // do nothing if context isn't initialized
	}
	l := v.(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()
	l.flush()
}

// logger lock must be held when calling this function.
func (l *logger) flush() {
	if l.flushed {
		return
	}
	for _, e := range l.entries {
		l.options.w.Write(l.options.format(e))
	}
	l.entries = nil // free up memory
	l.flushed = true
}

func log(ctx context.Context, sev Severity, buffer bool, msg string, keyvals ...interface{}) {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return // do nothing if context isn't initialized
	}
	l := v.(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()

	if !l.options.debug && sev == SeverityDebug {
		return
	}
	if l.options.debug && !l.flushed {
		l.flush()
	}

	keyvals = append(l.keyvals, keyvals...)
	keyvals = append(l.options.keyvals, keyvals...)
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, nil)
	}

	if len(msg) > l.options.maxsize {
		msg = msg[0:l.options.maxsize]
	}
	truncate(keyvals, l.options.maxsize)

	e := &Entry{timeNow().UTC(), sev, keyvals, msg}
	if l.flushed || !buffer {
		l.options.w.Write(l.options.format(e))
		return
	}
	l.entries = append(l.entries, e)
}

// Parse extracts the keys and values from the given key/value pairs. The
// resulting slices are of the same length and ordered in the same way.
func (kv KeyVals) Parse() (keys []string, vals []interface{}) {
	if len(kv) == 0 {
		return
	}
	keys = make([]string, len(kv)/2)
	vals = make([]interface{}, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			key = "<INVALID>"
		}
		keys[i/2] = key
		vals[i/2] = kv[i+1]
	}
	return keys, vals
}

// String returns a string representation of the log severity.
func (l Severity) String() string {
	switch l {
	case SeverityDebug:
		return "debug"
	case SeverityInfo:
		return "info"
	case SeverityError:
		return "error"
	default:
		return "<INVALID>"
	}
}

// Code returns a 4-character code for the log severity.
func (l Severity) Code() string {
	switch l {
	case SeverityDebug:
		return "DEBG"
	case SeverityInfo:
		return "INFO"
	case SeverityError:
		return "ERRO"
	default:
		return "<INVALID>"
	}
}

// Color returns an escape sequence that colors the output for the given
// severity.
func (l Severity) Color() string {
	switch l {
	case SeverityDebug:
		return "\033[37m"
	case SeverityInfo:
		return "\033[34m"
	case SeverityError:
		return "\033[1;31m"
	default:
		return ""
	}
}

// truncate makes sure that all string values in keyvals are no longer than
// maxsize and that all slice values are truncated to maxsize.
//
// Note: This could get very complicated very quickly (there could be different
// max values for strings and slices, it could compute total size for slices vs.
// size for each element, could recurse further etc.) - the point is to protect
// against obvious mistakes - not to implement a bullet-proof solution.
func truncate(keyvals []interface{}, maxsize int) {
	for i := 1; i < len(keyvals); i += 2 {
		switch v := keyvals[i].(type) {
		case string:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []string:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
			for j, s := range v {
				if len(s) > maxsize {
					v[j] = s[0:maxsize]
				}
			}
		case []int:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []int32:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []int64:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []uint:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []uint32:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []uint64:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []float32:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []float64:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []bool:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
		case []interface{}:
			if len(v) > maxsize {
				keyvals[i] = v[0:maxsize]
			}
			for j, e := range v {
				if s, ok := e.(string); ok && len(s) > maxsize {
					v[j] = s[0:maxsize]
				}
			}
		}
	}
}
