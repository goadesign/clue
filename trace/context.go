package trace

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type (
	// ctxKey is a private type used to store the tracer provider in the context.
	ctxKey int

	// stateBag tracks the provider, tracer and active span sequence for a request.
	stateBag struct {
		svc      string
		provider trace.TracerProvider
		tracer   trace.Tracer
		spans    []trace.Span
	}
)

const (
	// InstrumentationLibraryName is the name of the instrumentation library.
	InstrumentationLibraryName = "goa.design/clue"

	// AttributeRequestID is the name of the span attribute that contains the
	// request ID.
	AttributeRequestID = "request.id"
)

const (
	// stateKey is used to store the tracing state the context.
	stateKey ctxKey = iota + 1
)

// Context initializes the context so it can be used to create traces.
func Context(ctx context.Context, svc string, opts ...TraceOption) (context.Context, error) {
	options := defaultOptions()
	for _, o := range opts {
		err := o(ctx, options)
		if err != nil {
			return nil, err
		}
	}

	if options.disabled {
		return withProvider(ctx, trace.NewNoopTracerProvider(), svc), nil
	}

	if options.exporter == nil {
		return nil, errors.New("missing exporter")
	}

	res := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(svc))
	rootSampler := adaptiveSampler(options.maxSamplingRate, options.sampleSize)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(rootSampler)),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(options.exporter),
	)
	return withProvider(ctx, provider, svc), nil
}

// IsTraced returns true if the current request is traced.
func IsTraced(ctx context.Context) bool {
	span := trace.SpanFromContext(ctx)
	return span.IsRecording() && span.SpanContext().IsSampled()
}

// TraceProvider returns the underlying otel trace provider.
func TraceProvider(ctx context.Context) trace.TracerProvider {
	sb := ctx.Value(stateKey).(*stateBag)
	return sb.provider
}

// withProvider stores the tracer provider in the context.
func withProvider(ctx context.Context, provider trace.TracerProvider, svc string) context.Context {
	return context.WithValue(ctx, stateKey, &stateBag{provider: provider, svc: svc})
}

// withTracing initializes the tracing context, ctx must have been initialized
// with withProvider and the request must be traced by otel.
func withTracing(traceCtx, ctx context.Context) context.Context {
	state := traceCtx.Value(stateKey).(*stateBag)
	svc := state.svc
	provider := state.provider
	tracer := provider.Tracer(InstrumentationLibraryName)
	spans := []trace.Span{trace.SpanFromContext(ctx)}
	return context.WithValue(ctx, stateKey, &stateBag{svc, provider, tracer, spans})
}
