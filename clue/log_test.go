package clue

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"goa.design/clue/log"
)

func TestLog(t *testing.T) {
	traceID := trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8}
	spanID := trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8}
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "valid span context",
			ctx:  ctx,
		},
		{
			name: "invalid span context",
			ctx:  context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kvs := Log(tt.ctx)
			spanContext := trace.SpanFromContext(tt.ctx).SpanContext()

			if spanContext.IsValid() {
				require.Equal(t, 2, len(kvs))
				assert.Equal(t, log.TraceIDKey, kvs[0].K)
				assert.Equal(t, spanContext.TraceID().String(), kvs[0].V)
				assert.Equal(t, log.SpanIDKey, kvs[1].K)
				assert.Equal(t, spanContext.SpanID().String(), kvs[1].V)
			} else {
				assert.Equal(t, 0, len(kvs))
			}
		})
	}
}
