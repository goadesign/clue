package trace

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type (
	// ctxKey is a private type used to store the tracer provider in the context.
	ctxKey int

	// stateBag tracks the provider, tracer and active span sequence for a request.
	stateBag struct {
		svc        string
		provider   trace.TracerProvider
		propagator propagation.TextMapPropagator
		tracer     trace.Tracer
		spans      []trace.Span
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
		return withConfig(ctx, noop.NewTracerProvider(), options.propagator, svc), nil
	}

	if options.exporter == nil {
		return nil, errors.New("missing exporter, set one with the option 'trace.WithExporter'")
	}

	res := options.resource
	if res == nil {
		res = resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(svc))
	}

	rootSampler := adaptiveSampler(options.maxSamplingRate, options.sampleSize)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(rootSampler, options.parentSamplerOptions...)),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(options.exporter),
	)
	return withConfig(ctx, provider, options.propagator, svc), nil
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

// withConfig stores the clue tracing config in the context.
func withConfig(ctx context.Context, provider trace.TracerProvider, propagator propagation.TextMapPropagator, svc string) context.Context {
	return context.WithValue(ctx, stateKey, &stateBag{provider: provider, propagator: propagator, svc: svc})
}

// withTracing initializes the tracing context, ctx must have been initialized
// with withProvider and the request must be traced by otel.
func withTracing(traceCtx, ctx context.Context) context.Context {
	state := traceCtx.Value(stateKey).(*stateBag)
	svc := state.svc
	provider := state.provider
	propagator := state.propagator
	tracer := provider.Tracer(InstrumentationLibraryName)
	spans := []trace.Span{trace.SpanFromContext(ctx)}
	return context.WithValue(ctx, stateKey, &stateBag{
		svc:        svc,
		provider:   provider,
		propagator: propagator,
		tracer:     tracer,
		spans:      spans,
	})
}
