package instrument

import "github.com/prometheus/client_golang/prometheus"

type (
	// Option is a function that configures the instrumentation.
	Option func(*options)

	// options contains the configuration for the instrumentation.
	options struct {
		// durationBuckets is the buckets for the request duration histogram.
		durationBuckets []float64
		// requestSizeBuckets is the buckets for the request size histogram.
		requestSizeBuckets []float64
		// responseSizeBuckets is the buckets for the response size histogram.
		responseSizeBuckets []float64
		// Prometheus registerer
		registerer prometheus.Registerer
		// Prometheus gatherer
		gatherer prometheus.Gatherer
	}
)

// defaultOptions returns a new options struct with default values.
func defaultOptions() *options {
	return &options{
		durationBuckets:     []float64{10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
		requestSizeBuckets:  []float64{10, 100, 500, 1000, 5000, 10000, 50000, 100000, 1000000, 10000000},
		responseSizeBuckets: []float64{10, 100, 500, 1000, 5000, 10000, 50000, 100000, 1000000, 10000000},
		registerer:          prometheus.DefaultRegisterer,
		gatherer:            prometheus.DefaultGatherer,
	}
}

// WithDurationBuckets returns an option that sets the duration buckets for the
// request duration histogram.
func WithDurationBuckets(buckets []float64) Option {
	return func(c *options) {
		c.durationBuckets = buckets
	}
}

// WithRequestSizeBuckets returns an option that sets the request size buckets
// for the request size histogram.
func WithRequestSizeBuckets(buckets []float64) Option {
	return func(c *options) {
		c.requestSizeBuckets = buckets
	}
}

// WithResponseSizeBuckets returns an option that sets the response size buckets
// for the response size histogram.
func WithResponseSizeBuckets(buckets []float64) Option {
	return func(c *options) {
		c.responseSizeBuckets = buckets
	}
}

// WithRegisterer returns an option that sets the prometheus registerer.
func WithRegisterer(registerer prometheus.Registerer) Option {
	return func(c *options) {
		c.registerer = registerer
	}
}
