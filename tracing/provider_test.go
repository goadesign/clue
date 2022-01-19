package tracing

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTraceProvider(t *testing.T) {
	// We don't want to test otel, keep it simple...
	exporter := tracetest.NewInMemoryExporter()
	tp, err := NewTracerProvider(context.Background(), "test", exporter)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "span")
	if !span.IsRecording() {
		t.Errorf("span is not recording")
	}
	if !span.SpanContext().IsValid() {
		t.Errorf("span context is not valid")
	}
	if !span.SpanContext().IsSampled() {
		t.Errorf("span context is not sampled")
	}
	span.End()
}

func TestNoExporter(t *testing.T) {
	_, err := NewTracerProvider(context.Background(), "test", nil)
	if err != ErrNoExporter {
		t.Errorf("got error: %v, expected %v", err, ErrNoExporter)
	}
}
