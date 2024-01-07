package clue

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"goa.design/clue/log"
)

// Log is a log key/value pair generator function that can be used to log trace
// and span IDs. Example:
//
//	ctx := log.Context(ctx, WithFunc(trace.Log))
//	log.Printf(ctx, "message")
//
//	Output: traceID=<trace-id> spanID=<span-id> message
func Log(ctx context.Context) (kvs []log.KV) {
	spanContext := trace.SpanFromContext(ctx).SpanContext()
	if spanContext.IsValid() {
		kvs = append(kvs, log.KV{K: log.TraceIDKey, V: spanContext.TraceID().String()})
		kvs = append(kvs, log.KV{K: log.SpanIDKey, V: spanContext.SpanID().String()})
	}
	return
}
