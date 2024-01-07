package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Log is a log key/value pair generator function that can be used to log trace
// and span IDs. Example:
//
//	ctx := log.Context(ctx, WithFunc(trace.Log))
//	log.Printf(ctx, "message")
//
//	Output: traceID=<trace-id> spanID=<span-id> message
func Span(ctx context.Context) (kvs []KV) {
	spanContext := trace.SpanFromContext(ctx).SpanContext()
	if spanContext.IsValid() {
		kvs = append(kvs, KV{K: TraceIDKey, V: spanContext.TraceID().String()})
		kvs = append(kvs, KV{K: SpanIDKey, V: spanContext.SpanID().String()})
	}
	return
}
