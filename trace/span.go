package trace

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan starts a new span with the given name and attributes and stores it
// in the returned context if the request is traced, does nothing otherwise.
// keyvals must be a list of alternating keys and values.
func StartSpan(ctx context.Context, name string, keyvals ...string) context.Context {
	s := ctx.Value(stateKey)
	if s == nil {
		return ctx
	}
	ctx, span := s.(*stateBag).tracer.Start(ctx, name, trace.WithAttributes(toKeyVal(keyvals)...))
	setActiveSpans(ctx, append(activeSpans(ctx), span))
	return ctx
}

// End ends the current span if any.
func EndSpan(ctx context.Context) {
	spans := activeSpans(ctx)
	if len(spans) == 0 {
		return
	}
	spans[len(spans)-1].End()
	setActiveSpans(ctx, spans[:len(spans)-1])
}

// SetSpanAttributes adds the given attributes to the current span if any.
// keyvals must be a list of alternating keys and values. It overwrites any
// existing attributes with the same key.
func SetSpanAttributes(ctx context.Context, keyvals ...string) {
	span := activeSpan(ctx)
	if span == nil {
		return
	}
	span.SetAttributes(toKeyVal(keyvals)...)
}

// AddEvent records an event with the given name and attributes in the current
// span if any.
func AddEvent(ctx context.Context, name string, keyvals ...string) {
	span := activeSpan(ctx)
	if span == nil {
		return
	}
	kvs := toKeyVal(keyvals)
	span.AddEvent(name, trace.WithAttributes(kvs...))
}

// Succeed sets the status of the current span to success if any.
func Succeed(ctx context.Context) {
	span := activeSpan(ctx)
	if span == nil {
		return
	}
	span.SetStatus(codes.Ok, "")
}

// Fail sets the status of the current span to failed and attaches the failure
// message.
func Fail(ctx context.Context, msg string) {
	span := activeSpan(ctx)
	if span == nil {
		return
	}
	span.SetStatus(codes.Error, msg)
}

// RecordError records err as an exception span event for the current span if
// any. An additional call to SetStatus is required if the Status of the Span
// should be set to Error, as this method does not change the Span status.
func RecordError(ctx context.Context, err error) {
	span := activeSpan(ctx)
	if span == nil {
		return
	}
	span.RecordError(err)
}

// TraceID returns the trace ID of the current span if any, empty string otherwise.
func TraceID(ctx context.Context) string {
	span := activeSpan(ctx)
	if span == nil {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

// SpanID returns the span ID of the current span if any, empty string otherwise.
func SpanID(ctx context.Context) string {
	span := activeSpan(ctx)
	if span == nil {
		return ""
	}
	return span.SpanContext().SpanID().String()
}

// activeSpans returns the active spans of the tracing state.
func activeSpans(ctx context.Context) []trace.Span {
	s := ctx.Value(stateKey)
	if s == nil {
		return nil
	}
	return s.(*stateBag).spans
}

// setActiveSpans updates the active spans of the tracing state.
func setActiveSpans(ctx context.Context, spans []trace.Span) {
	s := ctx.Value(stateKey)
	if s == nil {
		return
	}
	s.(*stateBag).spans = spans
}

func activeSpan(ctx context.Context) trace.Span {
	spans := activeSpans(ctx)
	if len(spans) == 0 {
		return nil
	}
	return spans[len(spans)-1]
}

func toKeyVal(kvs []string) []attribute.KeyValue {
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "")
	}
	keyvals := make([]attribute.KeyValue, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		keyvals[i/2] = attribute.String(kvs[i], kvs[i+1])
	}
	return keyvals
}
