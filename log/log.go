package log

import (
	"context"
	"sync"
)

type (
	// KeyVals represents a list of key/value pairs.
	KeyVals []interface{}

	// Log entry
	Entry struct {
		Level   Level
		KeyVals KeyVals
		Message string
	}

	// Logger implementation
	logger struct {
		options *options
		lock    sync.Mutex
		keyvals []interface{}
		entries []*Entry
		flushed bool
	}

	// Log level enum
	Level int

	// private type for context keys
	ctxKey int
)

const (
	LvlDebug Level = iota + 1
	LvlInfo
	LvlError
)

const (
	ctxLogger ctxKey = iota + 1
)

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

// Debug logs a debug message.
func Debug(ctx context.Context, msg string, keyvals ...interface{}) {
	log(ctx, LvlDebug, true, msg, keyvals...)
}

// Print logs an info message and ignores buffering.
func Print(ctx context.Context, msg string, keyvals ...interface{}) {
	log(ctx, LvlInfo, false, msg, keyvals...)
}

// Info logs an info message.
func Info(ctx context.Context, msg string, keyvals ...interface{}) {
	log(ctx, LvlInfo, true, msg, keyvals...)
}

// Error logs an error message.
func Error(ctx context.Context, msg string, keyvals ...interface{}) {
	Flush(ctx)
	log(ctx, LvlError, true, msg, keyvals...)
}

// With adds the given key/value pairs to the log context.
func With(ctx context.Context, keyvals ...interface{}) context.Context {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return ctx
	}
	l := v.(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, nil)
	}
	l.keyvals = append(l.keyvals, keyvals...)
	return ctx
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

func log(ctx context.Context, level Level, buffer bool, msg string, keyvals ...interface{}) {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return // do nothing if context isn't initialized
	}
	l := v.(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()

	if !l.options.debug && level == LvlDebug {
		return
	}
	if l.options.debug && !l.flushed {
		l.flush()
	}

	keyvals = append(l.keyvals, keyvals...)
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, nil)
	}

	e := &Entry{level, keyvals, msg}
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

// String returns a string representation of the log level.
func (l Level) String() string {
	switch l {
	case LvlDebug:
		return "DEBUG"
	case LvlInfo:
		return "INFO"
	case LvlError:
		return "ERROR"
	default:
		return "<INVALID>"
	}
}

// Color returns an escape sequence that colors the output for the given level.
func (l Level) Color() string {
	switch l {
	case LvlDebug:
		return "\033[1;32m"
	case LvlInfo:
		return "\033[1;34m"
	case LvlError:
		return "\033[1;31m"
	default:
		return ""
	}
}
