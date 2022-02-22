package trace

import (
	"context"
	"testing"
)

func TestLog(t *testing.T) {
	ctx, _ := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())
	ctx = StartSpan(ctx, "span")
	kvs := Log(ctx)
	if len(kvs) != 2 {
		t.Fatalf("got %d kvs, expected 2", len(kvs))
	}
	if kvs[0].K != TraceIDLogKey {
		t.Errorf("got kvs[0].K %q, expected %q", kvs[0].K, TraceIDLogKey)
	}
	if kvs[0].V != TraceID(ctx) {
		t.Errorf("got kvs[0].V %q, expected %q", kvs[0].V, TraceID(ctx))
	}
	if kvs[1].K != SpanIDLogKey {
		t.Errorf("got kvs[1].K %q, expected %q", kvs[1].K, SpanIDLogKey)
	}
	if kvs[1].V != SpanID(ctx) {
		t.Errorf("got kvs[1].V %q, expected %q", kvs[1].V, SpanID(ctx))
	}
}
