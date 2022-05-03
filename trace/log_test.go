package trace

import (
	"context"
	"testing"

	"goa.design/clue/log"
)

func TestLog(t *testing.T) {
	ctx, _ := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())
	ctx = StartSpan(ctx, "span")
	kvs := Log(ctx)
	if len(kvs) != 2 {
		t.Fatalf("got %d kvs, expected 2", len(kvs))
	}
	if kvs[0].K != log.TraceIDKey {
		t.Errorf("got kvs[0].K %q, expected %q", kvs[0].K, log.TraceIDKey)
	}
	if kvs[0].V != TraceID(ctx) {
		t.Errorf("got kvs[0].V %q, expected %q", kvs[0].V, TraceID(ctx))
	}
	if kvs[1].K != log.SpanIDKey {
		t.Errorf("got kvs[1].K %q, expected %q", kvs[1].K, log.SpanIDKey)
	}
	if kvs[1].V != SpanID(ctx) {
		t.Errorf("got kvs[1].V %q, expected %q", kvs[1].V, SpanID(ctx))
	}
}
