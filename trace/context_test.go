package trace

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func testContext(provider trace.TracerProvider) context.Context {
	return withProvider(context.Background(), provider, propagation.TraceContext{}, "test")
}

func TestContext(t *testing.T) {
	// We don't want to test otel, keep it simple...
	exporter := tracetest.NewInMemoryExporter()
	ctx, err := Context(context.Background(), "test",
		WithMaxSamplingRate(3), WithSampleSize(20), WithExporter(exporter))
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
}

func TestDisabled(t *testing.T) {
	ctx, err := Context(context.Background(), "test", WithDisabled())
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
	_, span := st.provider.Tracer("test").Start(ctx, "test")
	if span.IsRecording() {
		t.Error("expected span to be disabled")
	}
}

func TestIsTraced(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	ctx, err := Context(context.Background(), "test",
		WithMaxSamplingRate(3), WithSampleSize(20), WithExporter(exporter))
	if err != nil {
		t.Fatal(err)
	}
	if IsTraced(ctx) {
		t.Error("expected not traced")
	}
	ctx, err = Context(ctx, "test", WithMaxSamplingRate(3), WithSampleSize(20), WithExporter(exporter))
	if err != nil {
		t.Fatal(err)
	}
	ctx = withTracing(ctx, context.Background())
	ctx = StartSpan(ctx, "test")
	if !IsTraced(ctx) {
		t.Error("expected traced")
	}
	EndSpan(ctx)
}
