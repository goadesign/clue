package trace

import (
	"context"

	"goa.design/clue/log"
)

var (
	// TraceIDLogKey is the key used to log the trace ID.
	TraceIDLogKey = "traceID"

	// SpanIDLogKey is the key used to log the span ID.
	SpanIDLogKey = "spanID"
)

// Log is a log key/value pair generator function that can be used to log trace
// and span IDs. Example:
//
//    ctx := log.Context(ctx, WithFunc(trace.Log))
//    log.Printf(ctx, "message")
//
//    Output: traceID=<trace-id> spanID=<span-id> message
func Log(ctx context.Context) (kvs []log.KV) {
	if id := TraceID(ctx); id != "" {
		kvs = append(kvs, log.KV{K: TraceIDLogKey, V: id})
	}
	if id := SpanID(ctx); id != "" {
		kvs = append(kvs, log.KV{K: SpanIDLogKey, V: id})
	}
	return
}
