package instrument

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	// handlerOption is a function that configures the handler.
	handlerOption func(*handlerOptions)

	handlerOptions struct {
		// registerer is the prometheus registerer.
		registerer prometheus.Registerer
		// gatherer is the prometheus gatherer.
		gatherer prometheus.Gatherer
	}
)

// Handler returns a HTTP handler that collect metrics and serves them using the
// Prometheus export formats. It uses the context logger configured via
// micro/log if any to log errors. By default Handler uses the default
// prometheus registry to gather metrics and to register its own metrics. Use
// options WithGatherer and WithHandlerRegisterer to override the default values.
func Handler(ctx context.Context, opts ...handlerOption) http.Handler {
	options := defaultHandlerOptions()
	for _, o := range opts {
		o(options)
	}
	return promhttp.InstrumentMetricHandler(options.registerer, promhttp.HandlerFor(options.gatherer,
		promhttp.HandlerOpts{
			ErrorLog: logger{ctx},
			Registry: options.registerer,
		}))
}

// WithHandlerRegisterer returns an option that sets the prometheus registerer.
func WithHandlerRegisterer(registerer prometheus.Registerer) handlerOption {
	return func(c *handlerOptions) {
		c.registerer = registerer
	}
}

// WithGatherer returns an option that sets the prometheus gatherer used to
// collect the metrics.
func WithGatherer(gatherer prometheus.Gatherer) handlerOption {
	return func(c *handlerOptions) {
		c.gatherer = gatherer
	}
}

// defaultHandlerOptions returns a new HandlerOption struct with default values.
func defaultHandlerOptions() *handlerOptions {
	return &handlerOptions{
		registerer: prometheus.DefaultRegisterer,
		gatherer:   prometheus.DefaultGatherer,
	}
}
