package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type (
	// ctxKey is a private type used to store the tracer provider in the context.
	ctxKey int

	// stateBag tracks the provider, tracer and active span sequence for a request.
	stateBag struct {
		svc      string
		provider *sdktrace.TracerProvider
		tracer   trace.Tracer
		spans    []trace.Span
	}
)

const (
	// InstrumentationLibraryName is the name of the instrumentation library.
	InstrumentationLibraryName = "goa.design/micro"

	// AttributeRequestID is the name of the span attribute that contains the
	// request ID.
	AttributeRequestID = "request.id"
)

const (
	// stateKey is used to store the tracing state the context.
	stateKey ctxKey = iota + 1
)

// Context initializes the context so it can be used to create traces.
func Context(ctx context.Context, svc string, conn *grpc.ClientConn, opts ...TraceOption) (context.Context, error) {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	res := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(svc))
	rootSampler := adaptiveSampler(options.maxSamplingRate, options.sampleSize)
	if options.exporter == nil {
		exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, err
		}
		options.exporter = exporter
	}
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

// withProvider stores the tracer provider in the context.
func withProvider(ctx context.Context, provider *sdktrace.TracerProvider, svc string) context.Context {
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
