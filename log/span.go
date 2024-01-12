package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Span is a log key/value pair generator function that can be used to log trace
// and span IDs. Usage:
//
//	ctx := log.Context(ctx, WithFunc(log.Span))
//	log.Printf(ctx, "message")
//
//	Output: trace_id=<trace id> span_id=<span id> message
func Span(ctx context.Context) (kvs []KV) {
	spanContext := trace.SpanFromContext(ctx).SpanContext()
	if spanContext.IsValid() {
		kvs = append(kvs, KV{K: TraceIDKey, V: spanContext.TraceID().String()})
		kvs = append(kvs, KV{K: SpanIDKey, V: spanContext.SpanID().String()})
	}
	return
}
