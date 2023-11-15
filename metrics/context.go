package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	// stateBag is stored in the context and used to initialize the
	// appropriate metrics for this package HTTP middleware and gRPC
	// interceptors. This state is only needed during initialization and is
	// not intended to be kept in request contexts.
	stateBag struct {
		options     *options
		svc         string
		httpMetrics *httpMetrics
		grpcMetrics *grpcMetrics
	}

	// httpMetrics is the set of HTTP Metrics used by this package interceptors.
	httpMetrics struct {
		// Durations is a histogram of the duration of requests.
		Durations *prometheus.HistogramVec
		// RequestSizes is a histogram of the size of requests.
		RequestSizes *prometheus.HistogramVec
		// ResponseSizes is a histogram of the size of responses.
		ResponseSizes *prometheus.HistogramVec
		// ActiveRequests is a gauge of the number of active requests.
		ActiveRequests *prometheus.GaugeVec
	}

	// grpcMetrics is the set of gRPC Metrics used by this package interceptors.
	grpcMetrics struct {
		// Durations is a histogram of the duration of requests.
		Durations *prometheus.HistogramVec
		// RequestSizes is a histogram of the size of requests.
		RequestSizes *prometheus.HistogramVec
		// ResponseSizes is a histogram of the size of responses.
		ResponseSizes *prometheus.HistogramVec
		// ActiveRequests is a gauge of the number of active requests.
		ActiveRequests *prometheus.GaugeVec
		// StreamMessageSizes is a histogram of the size of messages sent on the stream.
		StreamMessageSizes *prometheus.HistogramVec
		// StreamResultSizes is a histogram of the size of results sent on the stream.
		StreamResultSizes *prometheus.HistogramVec
	}

	// Private type used to define context keys.
	ctxKey int
)

const (
	// metricHTTPDuration is the name of the HTTP duration metric.
	metricHTTPDuration = "http_server_duration_ms"
	// metricHTTPRequestSize is the name of the HTTP request size metric.
	metricHTTPRequestSize = "http_server_request_size_bytes"
	// metricHTTPResponseSize is the name of the HTTP response size metric.
	metricHTTPResponseSize = "http_server_response_size_bytes"
	// metricHTTPActiveRequests is the name of the HTTP active requests metric.
	metricHTTPActiveRequests = "http_server_active_requests"
	// metricRPCDuration is the name of the gRPC request duration metric.
	metricRPCDuration = "rpc_server_duration_ms"
	// metricRPCActiveRequests is the name of the gRPC active requests metric.
	metricRPCActiveRequests = "rpc_server_active_requests"
	// metricRPCRequestSize is the name of the gRPC request size metric.
	metricRPCRequestSize = "rpc_server_request_size_bytes"
	// metricRPCResponseSize is the name of the gRPC response size metric.
	metricRPCResponseSize = "rpc_server_response_size_bytes"
	// metricRPCStreamMessageSize is the name of the gRPC stream message size metric.
	metricRPCStreamMessageSize = "rpc_server_stream_message_size_bytes"
	// metricRPCStreamResponseSize is the name of the gRPC stream response size metric.
	metricRPCStreamResponseSize = "rpc_server_stream_response_size_bytes"
	// labelGoaService is the name of the label containing the Goa service name.
	labelGoaService = "goa_service"
	// labelHTTPVerb is the name of the label containing the HTTP verb.
	labelHTTPVerb = "http_verb"
	// labelHTTPHost is the name of the label containing the HTTP host.
	labelHTTPHost = "http_host"
	// labelHTTPPath is the name of the label containing the HTTP URL path.
	labelHTTPPath = "http_path"
	// labelHTTPStatusCode is the name of the label containing the HTTP status code.
	labelHTTPStatusCode = "http_status_code"
	// labelRPCService is the name of the RPC service label.
	labelRPCService = "rpc_service"
	// labelRPCMethod is the name of the RPC method label.
	labelRPCMethod = "rpc_method"
	// labelRPCStatusCode is the name of the RPC status code label.
	labelRPCStatusCode = "rpc_status_code"
)

const (
	// Context key used to capture request length.
	ctxReqLen ctxKey = iota + 1
	// Context key used to store initialization state bag.
	stateBagKey
)

var (
	// httpLabels is the set of dynamic labels used for all metrics but
	// MetricHTTPActiveRequests.
	httpLabels = []string{labelHTTPVerb, labelHTTPHost, labelHTTPPath, labelHTTPStatusCode}

	// httpActiveRequestsLabels is the set of dynamic labels used for the
	// MetricHTTPActiveRequests metric.
	httpActiveRequestsLabels = []string{labelHTTPVerb, labelHTTPHost, labelHTTPPath}

	// rpcLabels is the default set of dynamic metric labels
	rpcLabels = []string{labelRPCService, labelRPCMethod, labelRPCStatusCode}

	// NoCode is the set of dynamic labels used for active gRPC requests
	// metric and stream message and result size metrics.
	rpcNoCodeLabels = []string{labelRPCService, labelRPCMethod}
)

// Context initializes the given context for the HTTP, UnaryInterceptor and
// StreamInterceptor functions.
func Context(ctx context.Context, svc string, opts ...Option) context.Context {
	if m := ctx.Value(stateBagKey); m != nil {
		return ctx
	}

	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	return context.WithValue(ctx, stateBagKey, &stateBag{options: options, svc: svc})
}

func (state *stateBag) HTTPMetrics() *httpMetrics {
	if state.httpMetrics != nil {
		return state.httpMetrics
	}

	durations := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricHTTPDuration,
		Help:        "Histogram of request durations in milliseconds.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
		Buckets:     state.options.durationBuckets,
	}, httpLabels)
	state.options.registerer.MustRegister(durations)

	reqSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricHTTPRequestSize,
		Help:        "Histogram of request sizes in bytes.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
		Buckets:     state.options.requestSizeBuckets,
	}, httpLabels)
	state.options.registerer.MustRegister(reqSizes)

	respSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricHTTPResponseSize,
		Help:        "Histogram of response sizes in bytes.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
		Buckets:     state.options.responseSizeBuckets,
	}, httpLabels)
	state.options.registerer.MustRegister(respSizes)

	activeReqs := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        metricHTTPActiveRequests,
		Help:        "Gauge of active requests.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
	}, httpActiveRequestsLabels)
	state.options.registerer.MustRegister(activeReqs)

	state.httpMetrics = &httpMetrics{
		Durations:      durations,
		RequestSizes:   reqSizes,
		ResponseSizes:  respSizes,
		ActiveRequests: activeReqs,
	}

	return state.httpMetrics
}

func (state *stateBag) GRPCMetrics() *grpcMetrics {
	if state.grpcMetrics != nil {
		return state.grpcMetrics
	}

	durations := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricRPCDuration,
		Help:        "Histogram of request durations in milliseconds.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
		Buckets:     state.options.durationBuckets,
	}, rpcLabels)
	state.options.registerer.MustRegister(durations)

	reqSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricRPCRequestSize,
		Help:        "Histogram of request sizes in bytes.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
		Buckets:     state.options.requestSizeBuckets,
	}, rpcLabels)
	state.options.registerer.MustRegister(reqSizes)

	respSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricRPCResponseSize,
		Help:        "Histogram of response sizes in bytes.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
		Buckets:     state.options.responseSizeBuckets,
	}, rpcLabels)
	state.options.registerer.MustRegister(respSizes)

	streamMsgSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricRPCStreamMessageSize,
		Help:        "Histogram of stream message sizes in bytes.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
		Buckets:     state.options.requestSizeBuckets,
	}, rpcNoCodeLabels)
	state.options.registerer.MustRegister(streamMsgSizes)

	streamResSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricRPCStreamResponseSize,
		Help:        "Histogram of stream response sizes in bytes.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
		Buckets:     state.options.responseSizeBuckets,
	}, rpcNoCodeLabels)
	state.options.registerer.MustRegister(streamResSizes)

	activeReqs := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        metricRPCActiveRequests,
		Help:        "Gauge of active requests.",
		ConstLabels: prometheus.Labels{labelGoaService: state.svc},
	}, rpcNoCodeLabels)
	state.options.registerer.MustRegister(activeReqs)

	state.grpcMetrics = &grpcMetrics{
		Durations:          durations,
		RequestSizes:       reqSizes,
		ResponseSizes:      respSizes,
		ActiveRequests:     activeReqs,
		StreamMessageSizes: streamMsgSizes,
		StreamResultSizes:  streamResSizes,
	}

	return state.grpcMetrics
}
