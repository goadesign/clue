package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

type (
	options struct {
		maxSamplingRate int
		sampleSize      int
		exporter        sdktrace.SpanExporter
		propagator      propagation.TextMapPropagator
		disabled        bool
	}

	// TraceOption is a function that configures a provider.
	TraceOption func(ctx context.Context, opts *options) error
)

// defaultOptions returns the default sampler options.
func defaultOptions() *options {
	return &options{
		maxSamplingRate: 2,
		sampleSize:      10,
		propagator:      propagation.TraceContext{},
	}
}

// WithMaxSamplingRate sets the maximum sampling rate in requests per second.
func WithMaxSamplingRate(rate int) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.maxSamplingRate = rate
		return nil
	}
}

// WithSampleSize sets the number of requests between two adjustments of the
// sampling rate.
func WithSampleSize(size int) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.sampleSize = size
		return nil
	}
}

// WithDisabled disables tracing, not for use in production.
func WithDisabled() TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.disabled = true
		return nil
	}
}

// WithExporter sets the exporter to use.
func WithExporter(exporter sdktrace.SpanExporter) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.exporter = exporter
		return nil
	}
}

// WithPropagator sets the otel propagators
func WithPropagator(propagator propagation.TextMapPropagator) TraceOption {
	return func(ctx context.Context, opts *options) error {
		opts.propagator = propagator
		return nil
	}
}

func WithGRPCExporter(conn *grpc.ClientConn) TraceOption {
	return func(ctx context.Context, opts *options) error {
		exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return err
		}
		opts.exporter = exporter
		return nil
	}
}
