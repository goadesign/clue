package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Option is a function that configures the metricsation.
	Option func(*options)

	// RouteResolver is a function that resolves the route of a request used
	// to label metrics.  Using a route resolver makes it possible to label
	// all routes matching a pattern with the same label. As an example
	// services using the github.com/go-chi/chi/v5 muxer can use
	// chi.RouteContext(r.Context()).RoutePattern().
	RouteResolver func(r *http.Request) string

	// options contains the configuration for the metricsation.
	options struct {
		// durationBuckets is the buckets for the request duration histogram.
		durationBuckets []float64
		// requestSizeBuckets is the buckets for the request size histogram.
		requestSizeBuckets []float64
		// responseSizeBuckets is the buckets for the response size histogram.
		responseSizeBuckets []float64
		// Prometheus registerer
		registerer prometheus.Registerer
		// RouteResolver is used to label metrics.
		resolver RouteResolver
	}
)

var (
	DefaultDurationBuckets     = []float64{10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}
	DefaultRequestSizeBuckets  = []float64{10, 100, 500, 1000, 5000, 10000, 50000, 100000, 1000000, 10000000}
	DefaultResponseSizeBuckets = []float64{10, 100, 500, 1000, 5000, 10000, 50000, 100000, 1000000, 10000000}
)

// defaultOptions returns a new options struct with default values.
func defaultOptions() *options {
	return &options{
		durationBuckets:     DefaultDurationBuckets,
		requestSizeBuckets:  DefaultRequestSizeBuckets,
		responseSizeBuckets: DefaultResponseSizeBuckets,
		registerer:          prometheus.DefaultRegisterer,
	}
}

// WithRouteResolver returns an option that sets the route resolver used to
// label metrics.  The default uses the request path.
func WithRouteResolver(resolver RouteResolver) Option {
	return func(o *options) {
		o.resolver = resolver
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
