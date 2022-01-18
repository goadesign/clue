package tracing

import (
	"context"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

func init() {
	// Soooo you can only call once or you get a panic.
	otel.SetTextMapPropagator(xray.Propagator{})
}

// NewTracerProvider returns a tracer provider that uses an adaptive sampler.
// svc is the service name and collectorAddr is the address the ADOT collector
// sidecar is listening on. The tracer provider is configured to use AWS X-Ray.
//
// See https://aws-otel.github.io/docs/getting-started/go-sdk/trace-manual-instr
// for more information.
func NewTracerProvider(ctx context.Context, svc, collectorAddr string, opts ...samplerOption) (*sdktrace.TracerProvider, error) {
	options := newDefaultOptions()
	for _, o := range opts {
		o(options)
	}

	// Create and start new OTLP trace exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(collectorAddr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		return nil, err
	}

	idg := xray.NewIDGenerator()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		// the service name used to display traces in backends
		semconv.ServiceNameKey.String(svc),
	)

	rootSampler := AdaptiveSampler(options.maxSamplingRate, options.sampleSize)
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(rootSampler)),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithIDGenerator(idg),
	), nil
}
