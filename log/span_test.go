package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestSpan(t *testing.T) {
	// Create a mock span context
	traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	})

	// Create a context with the mock span context
	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

	// Call the Span function
	kvs := Span(ctx)

	// Assert that the expected key-value pairs are returned
	assert.Equal(t, []KV{
		{K: TraceIDKey, V: "0102030405060708090a0b0c0d0e0f10"},
		{K: SpanIDKey, V: "0102030405060708"},
	}, kvs)
}
