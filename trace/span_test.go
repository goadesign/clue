package trace

import (
	"context"
	"errors"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func newTestTracingContext() (context.Context, *tracetest.InMemoryExporter) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	ctx := withProvider(context.Background(), provider)
	return ctx, exporter
}

func TestStartEndSpan(t *testing.T) {
	// Make sure StartSpan and EndSpan do not panic
	StartSpan(context.Background(), "noop")
	EndSpan(context.Background())

	// Setup
	ctx, exporter := newTestTracingContext()

	// Create span
	ctx = withTracing(ctx, context.Background())
	ctx = StartSpan(ctx, "span")
	if !trace.SpanFromContext(ctx).SpanContext().IsValid() {
		t.Error("expected valid span")
	}
	if !trace.SpanFromContext(ctx).IsRecording() {
		t.Error("expected recording span")
	}
	if !trace.SpanFromContext(ctx).SpanContext().IsSampled() {
		t.Error("expected sampled span")
	}
	EndSpan(ctx)

	// Verify
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, expected 1", len(spans))
	}
	if spans[0].Name != "span" {
		t.Errorf("got span name %q, expected %q", spans[0].Name, "span")
	}
	var zeroID trace.TraceID
	if spans[0].Parent.TraceID() != zeroID {
		t.Errorf("got parent traceID %s, expected %s", spans[0].Parent.TraceID(), zeroID)
	}
}

func TestChildSpan(t *testing.T) {
	// Setup
	ctx, exporter := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())

	// Create spans
	ctx = StartSpan(ctx, "parent")
	ctx = StartSpan(ctx, "child")
	EndSpan(ctx)
	EndSpan(ctx)

	// Verify )
	spans := exporter.GetSpans() // returns spans in reverse order
	if len(spans) != 2 {
		t.Fatalf("got %d spans, expected 2", len(spans))
	}
	if spans[1].Name != "parent" {
		t.Errorf("got span name %q, expected %q", spans[1].Name, "parent")
	}
	var zeroID trace.TraceID
	if spans[1].Parent.TraceID() != zeroID {
		t.Errorf("got parent span %v, expected %v", spans[1].Parent.TraceID(), zeroID)
	}
	if spans[0].Name != "child" {
		t.Errorf("got span name %q, expected %q", spans[0].Name, "child")
	}
	if spans[0].Parent.TraceID() != spans[1].SpanContext.TraceID() {
		t.Errorf("got parent span %v, expected %v", spans[0].Parent, spans[1])
	}
}

func TestSetSpanAttributesNoContext(t *testing.T) {
	SetSpanAttributes(context.Background(), "key", "value")
}

func TestSetSpanAttributes(t *testing.T) {
	cases := []struct {
		name     string
		keyvals  []string
		expected map[string]string
	}{
		{"nil", nil, nil},
		{"empty", []string{}, nil},
		{"one", []string{"key", "value"}, map[string]string{"key": "value"}},
		{"invalid", []string{"key", "value", "key2"}, map[string]string{"key": "value", "key2": ""}},
		{"override", []string{"key", "value", "key", "value2"}, map[string]string{"key": "value2"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Setup
			ctx, exporter := newTestTracingContext()
			ctx = withTracing(ctx, context.Background())

			// Create span and set attributes
			ctx = StartSpan(ctx, "span")
			SetSpanAttributes(ctx, c.keyvals...)
			EndSpan(ctx)

			// Verify
			spans := exporter.GetSpans()
			if len(spans) != 1 {
				t.Fatalf("got %d spans, expected 1", len(spans))
			}
			if len(spans[0].Attributes) != len(c.expected) {
				t.Errorf("got %d attributes, expected %d", len(spans[0].Attributes), len(c.expected))
			}
			for k, v := range c.expected {
				found := false
				for _, att := range spans[0].Attributes {
					if att.Key == attribute.Key(k) && att.Value.AsString() == v {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("attribute %s not found", k)
				}
			}
		})
	}
}

func TestAddEventNoContext(t *testing.T) {
	AddEvent(context.Background(), "event")
}

func TestAddEvent(t *testing.T) {
	// Setup
	ctx, exporter := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())

	// Create span and add event
	ctx = StartSpan(ctx, "span")
	AddEvent(ctx, "event", "key", "value")
	EndSpan(ctx)

	// Verify
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, expected 1", len(spans))
	}
	if len(spans[0].Events) != 1 {
		t.Errorf("got %d events, expected 1", len(spans[0].Events))
	}
	if spans[0].Events[0].Name != "event" {
		t.Errorf("got event name %q, expected %q", spans[0].Events[0].Name, "event")
	}
	if len(spans[0].Events[0].Attributes) != 1 {
		t.Fatalf("got %d attributes, expected 1", len(spans[0].Events[0].Attributes))
	}
	if spans[0].Events[0].Attributes[0].Key != "key" {
		t.Errorf("got event attribute key %q, expected %q", spans[0].Events[0].Attributes[0].Key, "key")
	}
	if spans[0].Events[0].Attributes[0].Value.AsString() != "value" {
		t.Errorf("got event attribute value %q, expected %q", spans[0].Events[0].Attributes[0].Value.AsString(), "value")
	}
}

func TestFailNoContext(t *testing.T) {
	Fail(context.Background(), "desc")
}

func TestFail(t *testing.T) {
	// Setup
	ctx, exporter := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())

	// Create span and set status
	ctx = StartSpan(ctx, "span")
	Fail(ctx, "desc")
	EndSpan(ctx)

	// Verify
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, expected 1", len(spans))
	}
	if spans[0].Status.Code != codes.Error {
		t.Errorf("got status code %d, expected %d", spans[0].Status.Code, codes.Error)
	}
	if spans[0].Status.Description != "desc" {
		t.Errorf("got status description %q, expected %q", spans[0].Status.Description, "desc")
	}
}

func TestSucceedNoContext(t *testing.T) {
	Succeed(context.Background())
}

func TestSucceed(t *testing.T) {
	// Setup
	ctx, exporter := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())

	// Create span and set status
	ctx = StartSpan(ctx, "span")
	Fail(ctx, "desc")
	Succeed(ctx)
	EndSpan(ctx)

	// Verify
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, expected 1", len(spans))
	}
	if spans[0].Status.Code != codes.Ok {
		t.Errorf("got status code %d, expected %d", spans[0].Status.Code, codes.Ok)
	}
}

func TestRecordErrorNoContext(t *testing.T) {
	RecordError(context.Background(), errors.New("err"))
}

func TestRecordError(t *testing.T) {
	// Setup
	ctx, exporter := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())

	// Create span and add event
	ctx = StartSpan(ctx, "span")
	recordedErr := errors.New("recorded error")
	RecordError(ctx, recordedErr)
	EndSpan(ctx)

	// Verify
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, expected 1", len(spans))
	}
	if len(spans[0].Events) != 1 {
		t.Errorf("got %d events, expected 1", len(spans[0].Events))
	}
	if spans[0].Events[0].Name != "exception" {
		t.Errorf("got event name %q, expected %q", spans[0].Events[0].Name, "exception")
	}
	if spans[0].Events[0].Attributes[0].Key != "exception.type" {
		t.Errorf("got event attribute key %q, expected %q", spans[0].Events[0].Attributes[0].Key, "exception.type")
	}
	if spans[0].Events[0].Attributes[0].Value.AsString() != "*errors.errorString" {
		t.Errorf("got event attribute value %q, expected %q", spans[0].Events[0].Attributes[0].Value.AsString(), "*errors.errorString")
	}
}

func TestTraceIDNoContext(t *testing.T) {
	tid := TraceID(context.Background())
	if tid != "" {
		t.Errorf("got trace ID %q, expected empty", tid)
	}
}

func TestTraceID(t *testing.T) {
	// Setup
	ctx, _ := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())

	// Create span and validate
	ctx = StartSpan(ctx, "span")
	if trace.SpanFromContext(ctx).SpanContext().TraceID().String() != TraceID(ctx) {
		t.Errorf("got trace ID %q, expected %q", trace.SpanFromContext(ctx).SpanContext().TraceID().String(), TraceID(ctx))
	}
	EndSpan(ctx)
}

func TestSpanIDNoContext(t *testing.T) {
	sid := SpanID(context.Background())
	if sid != "" {
		t.Errorf("got span ID %q, expected empty", sid)
	}
}

func TestSpanID(t *testing.T) {
	// Setup
	ctx, _ := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())

	// Create span and add event
	ctx = StartSpan(ctx, "span")
	if trace.SpanFromContext(ctx).SpanContext().SpanID().String() != SpanID(ctx) {
		t.Errorf("got span ID %q, expected %q", trace.SpanFromContext(ctx).SpanContext().SpanID().String(), SpanID(ctx))
	}
	EndSpan(ctx)
}

func TestSetActiveSpansNoContext(t *testing.T) {
	setActiveSpans(context.Background(), nil)
}

func TestActiveSpans(t *testing.T) {
	// Setup
	ctx, _ := newTestTracingContext()
	ctx = withTracing(ctx, context.Background())

	// Create out-of-state span
	ctx, span := ctx.Value(stateKey).(*stateBag).tracer.Start(ctx, "span")

	// Make sure it's not in state
	spans := activeSpans(ctx)
	if len(spans) != 1 {
		t.Fatalf("got %d active spans, expected 1", len(spans))
	}

	// Manually add it to state
	setActiveSpans(ctx, []trace.Span{span})
	spans = activeSpans(ctx)
	if len(spans) != 1 {
		t.Fatalf("got %d active spans, expected 1", len(spans))
	}
	if spans[0] != span {
		t.Errorf("got active span %#v, expected %#v", spans[0], span)
	}

	// End span
	EndSpan(ctx)

	// Make sure it's not in state anymore
	spans = activeSpans(ctx)
	if len(spans) != 0 {
		t.Fatalf("got %d active spans, expected 0", len(spans))
	}
}
