package tracing

import sdktrace "go.opentelemetry.io/otel/sdk/trace"

type (
	SpanExporter = sdktrace.SpanExporter

	samplerOption func(opts *options)

	options struct {
		maxSamplingRate int
		sampleSize      int
	}
)

func newDefaultOptions() *options {
	return &options{
		maxSamplingRate: 2,
		sampleSize:      10,
	}
}

// WithMaxSamplingRate sets the maximum sampling rate in requests per second.
func WithMaxSamplingRate(rate int) samplerOption {
	return func(opts *options) {
		opts.maxSamplingRate = rate
	}
}

// WithSampleSize sets the number of requests between two adjustments of the
// sampling rate.
func WithSampleSize(size int) samplerOption {
	return func(opts *options) {
		opts.sampleSize = size
	}
}
