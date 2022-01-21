package tracing

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestContext(t *testing.T) {
	// We don't want to test otel, keep it simple...
	exporter := tracetest.NewInMemoryExporter()
	ctx, err := Context(context.Background(), "test", nil,
		WithMaxSamplingRate(3), WithSampleSize(20), withExporter(exporter))
	if err != nil {
		t.Fatal(err)
	}
	s := ctx.Value(stateKey)
	if s == nil {
		t.Fatal("expected state in context")
	}
	st, ok := s.(*stateBag)
	if !ok {
		t.Fatalf("got %T, expected *stateBag", s)
	}
	if st.provider == nil {
		t.Error("expected provider in tracing context")
	}
	if st.tracer == nil {
		t.Error("expected tracer in tracing context")
	}
	if len(st.spans) != 0 {
		t.Errorf("expected no span, got %d", len(st.spans))
	}
}

func TestIsTraced(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	ctx, err := Context(context.Background(), "test", nil,
		WithMaxSamplingRate(3), WithSampleSize(20), withExporter(exporter))
	if err != nil {
		t.Fatal(err)
	}
	if IsTraced(ctx) {
		t.Error("expected not traced")
	}
	ctx, err = Context(ctx, "test", nil, WithMaxSamplingRate(3), WithSampleSize(20), withExporter(exporter))
	if err != nil {
		t.Fatal(err)
	}
	ctx = StartSpan(ctx, "test")
	if !IsTraced(ctx) {
		t.Error("expected traced")
	}
	EndSpan(ctx)
}
