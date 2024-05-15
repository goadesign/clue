package clue

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	// Option is a function that initializes the clue configuration.
	Option func(*options)

	// options contains the clue configuration options.
	options struct {
		// readerInterval is the interval at which the metrics reader is
		// invoked.
		readerInterval time.Duration
		// maxSamplingRate is the maximum sampling rate for the trace exporter.
		maxSamplingRate int
		// sampleSize is the number of requests between two adjustments of the
		// sampling rate.
		sampleSize int
		// propagators is the trace propagators.
		propagators propagation.TextMapPropagator
		// resource is the resource containing any additional attributes.
		resource *resource.Resource
		// errorHandler is the error handler used by the otel package.
		errorHandler otel.ErrorHandler
	}
)

// defaultOptions returns a new options struct with default values.
// The logger in ctx is used to log errors.
func defaultOptions(ctx context.Context) *options {
	return &options{
		maxSamplingRate: 2,
		sampleSize:      10,
		propagators:     propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
		errorHandler:    NewErrorHandler(ctx),
	}
}

// WithReaderInterval returns an option that sets the interval at which the
// metrics reader is invoked.
func WithReaderInterval(interval time.Duration) Option {
	return func(c *options) {
		c.readerInterval = interval
	}
}

// WithMaxSamplingRate sets the maximum sampling rate in requests per second.
func WithMaxSamplingRate(rate int) Option {
	return func(opts *options) {
		opts.maxSamplingRate = rate
	}
}

// WithSampleSize sets the number of requests between two adjustments of the
// sampling rate.
func WithSampleSize(size int) Option {
	return func(opts *options) {
		opts.sampleSize = size
	}
}

// WithPropagators sets the propagators when extracting and injecting trace
// context.
func WithPropagators(propagator propagation.TextMapPropagator) Option {
	return func(opts *options) {
		opts.propagators = propagator
	}
}

// WithResource sets the resource containing any additional attributes.
func WithResource(res *resource.Resource) Option {
	return func(opts *options) {
		opts.resource = res
	}
}

// WithErrorHandler sets the error handler used by the telemetry package.
func WithErrorHandler(errorHandler otel.ErrorHandler) Option {
	return func(opts *options) {
		opts.errorHandler = errorHandler
	}
}
