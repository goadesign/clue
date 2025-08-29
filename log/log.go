package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	otellog "go.opentelemetry.io/otel/log"
)

type (
	// Log entry
	Entry struct {
		Time     time.Time
		Severity Severity
		KeyVals  kvList
	}

	// Logger implementation
	logger struct {
		options     *options
		lock        sync.Mutex
		keyvals     kvList
		entries     []*Entry
		flushed     bool
		otellog     otellog.Logger // OpenTelemetry logger
		otelLogging bool           // Guard against recursive OTEL logging
	}

	// Log severity enum
	Severity int

	// private type for context keys
	ctxKey int
)

const (
	SeverityDebug Severity = iota + 1
	SeverityInfo
	SeverityWarn
	SeverityError
)

// Be kind to tests
var (
	timeNow   = time.Now
	timeSince = time.Since
	osExit    = os.Exit
)

// Export color codes
var (
	ColorSeverityDebug = "\033[37m"
	ColorSeverityInfo  = "\033[34m"
	ColorSeverityWarn  = "\033[33m"
	ColorSeverityError = "\033[1;31m"
)

// Debug writes the key/value pairs to the log output if the log context is
// configured to log debug messages (via WithDebug).
func Debug(ctx context.Context, keyvals ...Fielder) {
	log(ctx, SeverityDebug, true, keyvals)
}

// Debugf sets the key MessageKey (default "msg") and calls Debug. Arguments
// are handled in the manner of fmt.Printf.
func Debugf(ctx context.Context, format string, v ...any) {
	Debug(ctx, KV{MessageKey, fmt.Sprintf(format, v...)})
}

// Print writes the key/value pairs to the log output ignoring buffering.
func Print(ctx context.Context, keyvals ...Fielder) {
	log(ctx, SeverityInfo, false, keyvals)
}

// Printf sets the key MessageKey (default "msg") and calls Print. Arguments
// are handled in the manner of fmt.Printf.
func Printf(ctx context.Context, format string, v ...any) {
	Print(ctx, KV{MessageKey, fmt.Sprintf(format, v...)})
}

// Info writes the key/value pairs to the log buffer or output if buffering is
// disabled.
func Info(ctx context.Context, keyvals ...Fielder) {
	log(ctx, SeverityInfo, true, keyvals)
}

// Infof sets the key MessageKey (default "msg") and calls Info. Arguments are
// handled in the manner of fmt.Printf.
func Infof(ctx context.Context, format string, v ...any) {
	Info(ctx, KV{MessageKey, fmt.Sprintf(format, v...)})
}

// Warn writes the key/value pairs to the log buffer or output if buffering is
// disabled, with SeverityWarn.
func Warn(ctx context.Context, keyvals ...Fielder) {
	log(ctx, SeverityWarn, true, keyvals)
}

// Warnf sets the key MessageKey (default "msg") and calls Warn. Arguments are
// handled in the manner of fmt.Printf.
func Warnf(ctx context.Context, format string, v ...any) {
	Warn(ctx, KV{MessageKey, fmt.Sprintf(format, v...)})
}

// Error flushes the log buffer and disables buffering if not already disabled.
// Error then sets the ErrorMessageKey (default "err") key with the given error
// and writes the key/value pairs to the log output.
func Error(ctx context.Context, err error, keyvals ...Fielder) {
	FlushAndDisableBuffering(ctx)
	if err != nil {
		kvs := make([]Fielder, len(keyvals)+1)
		copy(kvs[1:], keyvals)
		kvs[0] = KV{ErrorMessageKey, err.Error()}
		keyvals = kvs
	}
	log(ctx, SeverityError, true, keyvals)
}

// Errorf sets the key MessageKey (default "msg") and calls Error. Arguments
// are handled in the manner of fmt.Printf.
func Errorf(ctx context.Context, err error, format string, v ...any) {
	Error(ctx, err, KV{MessageKey, fmt.Sprintf(format, v...)})
}

// Fatal is equivalent to Error followed by a call to os.Exit(1)
func Fatal(ctx context.Context, err error, keyvals ...Fielder) {
	Error(ctx, err, keyvals...)
	osExit(1)
}

// Fatalf is equivalent to Errorf followed by a call to os.Exit(1)
func Fatalf(ctx context.Context, err error, format string, v ...any) {
	Fatal(ctx, err, KV{MessageKey, fmt.Sprintf(format, v...)})
}

// With creates a copy of the given log context and appends the given key/value
// pairs to it. Values must be strings, numbers, booleans, nil or a slice of
// these types.
func With(ctx context.Context, keyvals ...Fielder) context.Context {
	v := ctx.Value(ctxLogger)
	if v == nil {
		return ctx
	}
	l := v.(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()
	newLogger := logger{
		options: l.options,
		entries: l.entries,
		keyvals: l.keyvals.merge(keyvals),
		flushed: l.flushed,
		otellog: l.otellog, // Copy OTEL logger
	}
	if l.options.disableBuffering != nil && l.options.disableBuffering(ctx) {
		l.flush()
		newLogger.flushed = true
	} else {
		newLogger.entries = make([]*Entry, len(l.entries))
		copy(newLogger.entries, l.entries)
	}

	return context.WithValue(ctx, ctxLogger, &newLogger)
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

// writeEntry writes the log entry to both the local writer and OpenTelemetry
func (l *logger) writeEntry(e *Entry) {
	// Write to local writer
	l.options.w.Write(l.options.format(e)) // nolint: errcheck

	if l.otellog == nil {
		return
	}

	// Convert severity to OTEL level
	var level otellog.Severity
	switch e.Severity {
	case SeverityDebug:
		level = otellog.SeverityDebug
	case SeverityInfo:
		level = otellog.SeverityInfo
	case SeverityWarn:
		level = otellog.SeverityWarn
	case SeverityError:
		level = otellog.SeverityError
	default:
		level = otellog.SeverityInfo
	}

	// Create OTEL record
	record := otellog.Record{}
	record.SetSeverity(level)
	record.SetTimestamp(e.Time)

	// Extract message and build structured body
	var message string
	var attributes []otellog.KeyValue
	var bodyParts []string

	for _, kv := range e.KeyVals {
		if kv.K == MessageKey || kv.K == ErrorMessageKey {
			message = fmt.Sprintf("%v", kv.V)
		} else {
			// Add to attributes with proper type handling
			switch v := kv.V.(type) {
			case string:
				attributes = append(attributes, otellog.String(kv.K, v))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%q", kv.K, v))
			case int:
				attributes = append(attributes, otellog.Int64(kv.K, int64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case int8:
				attributes = append(attributes, otellog.Int64(kv.K, int64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case int16:
				attributes = append(attributes, otellog.Int64(kv.K, int64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case int32:
				attributes = append(attributes, otellog.Int64(kv.K, int64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case int64:
				attributes = append(attributes, otellog.Int64(kv.K, v))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case uint:
				attributes = append(attributes, otellog.Int64(kv.K, int64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case uint8:
				attributes = append(attributes, otellog.Int64(kv.K, int64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case uint16:
				attributes = append(attributes, otellog.Int64(kv.K, int64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case uint32:
				attributes = append(attributes, otellog.Int64(kv.K, int64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case uint64:
				attributes = append(attributes, otellog.Int64(kv.K, int64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%d", kv.K, v))
			case float32:
				attributes = append(attributes, otellog.Float64(kv.K, float64(v)))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%g", kv.K, v))
			case float64:
				attributes = append(attributes, otellog.Float64(kv.K, v))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%g", kv.K, v))
			case bool:
				attributes = append(attributes, otellog.Bool(kv.K, v))
				bodyParts = append(bodyParts, fmt.Sprintf("%s=%t", kv.K, v))
			default:
				// Check for Stringer interface
				if stringer, ok := v.(fmt.Stringer); ok {
					attributes = append(attributes, otellog.String(kv.K, stringer.String()))
					bodyParts = append(bodyParts, fmt.Sprintf("%s=%q", kv.K, stringer.String()))
				} else if jsonMarshaler, ok := v.(json.Marshaler); ok {
					// Check for MarshalJSON interface
					if data, err := jsonMarshaler.MarshalJSON(); err == nil {
						// Remove quotes from JSON string if it's a string
						jsonStr := string(data)
						if len(jsonStr) >= 2 && jsonStr[0] == '"' && jsonStr[len(jsonStr)-1] == '"' {
							jsonStr = jsonStr[1 : len(jsonStr)-1]
						}
						attributes = append(attributes, otellog.String(kv.K, jsonStr))
						bodyParts = append(bodyParts, fmt.Sprintf("%s=%q", kv.K, jsonStr))
					} else {
						// Fallback to string representation if MarshalJSON fails
						attributes = append(attributes, otellog.String(kv.K, fmt.Sprintf("%v", v)))
						bodyParts = append(bodyParts, fmt.Sprintf("%s=%q", kv.K, v))
					}
				} else {
					// For any other type, convert to string
					attributes = append(attributes, otellog.String(kv.K, fmt.Sprintf("%v", v)))
					bodyParts = append(bodyParts, fmt.Sprintf("%s=%q", kv.K, v))
				}
			}
		}
	}

	// Build the structured body: "message key1=value1 key2=value2"
	var body string
	if message != "" {
		body = message
	}
	if len(bodyParts) > 0 {
		if body != "" {
			body += " "
		}
		body += strings.Join(bodyParts, " ")
	}

	// Set the body to the structured format
	record.SetBody(otellog.StringValue(body))

	// Add all keyvals as attributes for proper OTEL structure
	for _, attr := range attributes {
		record.AddAttributes(attr)
	}

	// Log to OTEL - use a background context to avoid recursive logging
	// This prevents the OTEL logger from calling back into the clue logger
	otelCtx := context.Background()

	// Add timeout to prevent hanging
	ctx, cancel := context.WithTimeout(otelCtx, 5*time.Second)
	defer cancel()

	// Use a goroutine to prevent blocking
	go func() {
		l.otellog.Emit(ctx, record)
	}()
}

func (l *logger) flush() {
	if l.flushed {
		return
	}
	for _, e := range l.entries {
		l.writeEntry(e)
	}
	l.entries = nil // free up memory
	l.flushed = true
}

func log(ctx context.Context, sev Severity, buffer bool, fielders []Fielder) {
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

	var keyvals kvList
	keyvals = keyvals.merge(fielders)
	keyvals = append(l.keyvals, keyvals...)
	keyvals = append(l.options.keyvals, keyvals...)
	for _, fn := range l.options.kvfuncs {
		keyvals = append(keyvals, fn(ctx)...)
	}
	truncate(keyvals, l.options.maxsize)

	e := &Entry{timeNow().UTC(), sev, keyvals}
	if l.flushed || !buffer {
		l.writeEntry(e)
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
	case SeverityWarn:
		return "warn"
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
	case SeverityWarn:
		return "WARN"
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
		return ColorSeverityDebug
	case SeverityInfo:
		return ColorSeverityInfo
	case SeverityWarn:
		return ColorSeverityWarn
	case SeverityError:
		return ColorSeverityError
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
		lw.Writer.Write(b) // nolint: errcheck
		return lw.max - lw.n, errTruncated
	}
	lw.n = newLen
	return lw.Writer.Write(b)
}
