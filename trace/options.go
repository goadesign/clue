package trace

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type (
	options struct {
		maxSamplingRate int
		sampleSize      int
		exporter        sdktrace.SpanExporter
		disabled        bool
	}

	// TraceOption is a function that configures a provider.
	TraceOption func(opts *options)
)

// defaultOptions returns the default sampler options.
func defaultOptions() *options {
	return &options{
		maxSamplingRate: 2,
		sampleSize:      10,
	}
}

// WithMaxSamplingRate sets the maximum sampling rate in requests per second.
func WithMaxSamplingRate(rate int) TraceOption {
	return func(opts *options) {
		opts.maxSamplingRate = rate
	}
}

// WithSampleSize sets the number of requests between two adjustments of the
// sampling rate.
func WithSampleSize(size int) TraceOption {
	return func(opts *options) {
		opts.sampleSize = size
	}
}

// WithDisabled disables tracing, not for use in production.
func WithDisabled() TraceOption {
	return func(opts *options) {
		opts.disabled = true
	}
}

// withExporter sets the exporter to use. This is intended for tests.
func withExporter(exporter sdktrace.SpanExporter) TraceOption {
	return func(opts *options) {
		opts.exporter = exporter
	}
}
