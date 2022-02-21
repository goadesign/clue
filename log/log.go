package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type (
	// KV represents a key/value pair. Values must be strings, numbers,
	// booleans, nil or a slice of these types.
	KV struct {
		K string
		V interface{}
	}

	// Log entry
	Entry struct {
		Time     time.Time
		Severity Severity
		KeyVals  []KV
	}

	// Logger implementation
	logger struct {
		options *options
		lock    sync.Mutex
		keyvals []KV
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
var (
	timeNow = time.Now
	osExit  = os.Exit
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

// Debug writes the key/value pairs to the log output if the log context is
// configured to log debug messages (via WithDebug).
func Debug(ctx context.Context, keyvals ...KV) {
	log(ctx, SeverityDebug, true, keyvals)
}

// Debugf sets the key "msg" and calls Debug. Arguments are handled in the
// manner of fmt.Printf.
func Debugf(ctx context.Context, format string, v ...interface{}) {
	Debug(ctx, KV{"msg", fmt.Sprintf(format, v...)})
}

// Print writes the key/value pairs to the log output ignoring buffering.
func Print(ctx context.Context, keyvals ...KV) {
	log(ctx, SeverityInfo, false, keyvals)
}

// Printf sets the key "msg" and calls Print. Arguments are handled in the
// manner of fmt.Printf.
func Printf(ctx context.Context, format string, v ...interface{}) {
	Print(ctx, KV{"msg", fmt.Sprintf(format, v...)})
}

// Info writes the key/value pairs to the log buffer or output if buffering is
// disabled.
func Info(ctx context.Context, keyvals ...KV) {
	log(ctx, SeverityInfo, true, keyvals)
}

// Infof sets the key "msg" and calls Info. Arguments are handled in the manner
// of fmt.Printf.
func Infof(ctx context.Context, format string, v ...interface{}) {
	Info(ctx, KV{"msg", fmt.Sprintf(format, v...)})
}

// Error flushes the log buffer and disables buffering if not already disabled.
// Error then sets the "err" key with the given error and writes the key/value
// pairs to the log output.
func Error(ctx context.Context, err error, keyvals ...KV) {
	FlushAndDisableBuffering(ctx)
	if err != nil {
		keyvals = append(keyvals, KV{"err", err.Error()})
	}
	log(ctx, SeverityError, true, keyvals)
}

// Errorf sets the key "msg" and calls Error. Arguments are handled in the
// manner of fmt.Printf.
func Errorf(ctx context.Context, err error, format string, v ...interface{}) {
	Error(ctx, err, KV{"msg", fmt.Sprintf(format, v...)})
}

// Fatal is equivalent to Error followed by a call to os.Exit(1)
func Fatal(ctx context.Context, err error, keyvals ...KV) {
	Error(ctx, err, keyvals...)
	osExit(1)
}

// Fatalf is equivalent to Errorf followed by a call to os.Exit(1)
func Fatalf(ctx context.Context, err error, format string, v ...interface{}) {
	Fatal(ctx, err, KV{"msg", fmt.Sprintf(format, v...)})
}

// With creates a copy of the given log context and appends the given key/value
// pairs to it. Values must be strings, numbers, booleans, nil or a slice of
// these types.
func With(ctx context.Context, keyvals ...KV) context.Context {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return ctx
	}
	l := v.(*logger)
	l.lock.Lock()
	copy := logger{
		options: l.options,
		entries: l.entries,
		keyvals: append(l.keyvals, keyvals...),
		flushed: l.flushed,
	}
	l.lock.Unlock()

	// Make sure that if Go needs to grow the slice then each context gets
	// its own memory.
	copy.entries = copy.entries[:len(copy.entries):len(copy.entries)]

	return context.WithValue(ctx, ctxLogger, &copy)
}

// FlushAndDisableBuffering flushes the log entries to the writer and stops
// buffering the given context.
func FlushAndDisableBuffering(ctx context.Context) {
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

func log(ctx context.Context, sev Severity, buffer bool, keyvals []KV) {
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
	for _, fn := range l.options.kvfuncs {
		keyvals = append(keyvals, fn(ctx)...)
	}
	truncate(keyvals, l.options.maxsize)

	e := &Entry{timeNow().UTC(), sev, keyvals}
	if l.flushed || !buffer {
		l.options.w.Write(l.options.format(e))
		return
	}
	l.entries = append(l.entries, e)
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

const truncationSuffix = " ... <clue/log.truncated>"

var errTruncated = errors.New("truncated value")

// truncate makes sure that all string values in keyvals are no longer than
// maxsize and that all slice values are truncated to maxsize.
//
// Note: This could get very complicated very quickly (there could be different
// max values for strings and slices, it could compute total size for slices vs.
// size for each element, could recurse further etc.) - the point is to protect
// against obvious mistakes - not to implement a bullet-proof solution.
func truncate(keyvals []KV, maxsize int) {
	if len(keyvals) > maxsize {
		keyvals = keyvals[:maxsize]
		keyvals = append(keyvals, KV{"log", truncationSuffix})
	}
	for i, kv := range keyvals {
		switch kv.V.(type) {
		case int, int8, int16, int32, int64:
			continue
		case uint, uint8, uint16, uint32, uint64:
			continue
		case float32, float64:
			continue
		case bool:
			continue
		case nil:
			continue
		default:
			var buf bytes.Buffer
			_, err := fmt.Fprintf(newLimitWriter(&buf, maxsize), "%v", kv.V)
			if errors.Is(err, errTruncated) {
				fmt.Fprint(&buf, truncationSuffix)
				keyvals[i] = KV{K: kv.K, V: buf.String()}
			}
		}
	}
}

type limitWriter struct {
	io.Writer
	max int
	n   int
}

func newLimitWriter(w io.Writer, max int) io.Writer {
	return &limitWriter{
		Writer: w,
		max:    max,
	}
}

func (lw *limitWriter) Write(b []byte) (int, error) {
	newLen := lw.n + len(b)
	if newLen > lw.max {
		b = b[:lw.max-lw.n]
		lw.Writer.Write(b)
		return lw.max - lw.n, errTruncated
	}
	lw.n = newLen
	return lw.Writer.Write(b)
}
