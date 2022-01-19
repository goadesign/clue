package tracing

import (
	"context"
	"errors"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// TracerProvider creates tracers that are configured to export traces to AWS
// X-Ray.
type TracerProvider = trace.TracerProvider

// AttributeRequestID is the name of the span attribute that contains the
// request ID.
const AttributeRequestID = "request_id"

// ErrNoExporter is returned when no exporter is configured.
var ErrNoExporter = errors.New("no exporter configured")

// NewTracerProvider returns a tracer provider that uses an adaptive sampler.
// svc is the service name and collectorAddr is the address the ADOT collector
// sidecar is listening on. The tracer provider is configured to use AWS X-Ray.
//
// See https://aws-otel.github.io/docs/getting-started/go-sdk/trace-manual-instr
// for more information.
func NewTracerProvider(ctx context.Context, svc string, exporter SpanExporter, opts ...samplerOption) (TracerProvider, error) {
	if exporter == nil {
		return nil, ErrNoExporter
	}
	options := newDefaultOptions()
	for _, o := range opts {
		o(options)
	}

	idg := xray.NewIDGenerator()
	res := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(svc))
	rootSampler := AdaptiveSampler(options.maxSamplingRate, options.sampleSize)

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(rootSampler)),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithIDGenerator(idg),
	), nil
}

// NewGRPCSpanExporter returns a new span exporter that uses the given host
// address to send spans via gRPC to the remote OLTP span collector.
func NewGRPCSpanExporter(ctx context.Context, addr string) (SpanExporter, error) {
	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(addr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
}
