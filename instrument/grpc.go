package instrument

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type (
	// streamWrapper wraps a grpc.ServerStream with prometheus metrics.
	streamWrapper struct {
		grpc.ServerStream
		labels              prometheus.Labels
		reqSizes, respSizes *prometheus.HistogramVec
	}
)

const (
	// MetricRPCDuration is the name of the gRPC request duration metric.
	MetricRPCDuration = "rpc_server_duration_ms"
	// MetricRPCActiveRequests is the name of the gRPC active requests metric_
	MetricRPCActiveRequests = "rpc_server_active_requests"
	// MetricRPCRequestSize is the name of the gRPC request size metric_
	MetricRPCRequestSize = "rpc_server_request_size_bytes"
	// MetricRPCResponseSize is the name of the gRPC response size metric_
	MetricRPCResponseSize = "rpc_server_response_size_bytes"
	// LabelPeerAddr is the peer host address_
	LabelPeerAddr = "net_peer_addr"
	// LabelRPCMethod is the name of the RPC method label_
	LabelRPCMethod = "rpc_method"
	// LabelRPCStatusCode is the name of the RPC status code label_
	LabelRPCStatusCode = "rpc_status_code"
)

var (
	// RPCLabels is the set of dynamic labels used for all metrics but
	// MetricRPCActiveRequests.
	RPCLabels = []string{LabelGoaMethod, LabelPeerAddr, LabelRPCMethod, LabelRPCStatusCode}

	// RPCActiveRequestsLabels is the set of dynamic labels used for
	// MetricRPCActiveRequests.
	RPCActiveRequestsLabels = []string{LabelRPCMethod, LabelPeerAddr}
)

// UnaryServerInterceptor creates a gRPC unary server interceptor that instruments the
// requests. The returned interceptor adds the following metrics:
//
//    * `grpc.server.duration`: Histogram of request durations in milliseconds.
//    * `grpc.server.active_requests`: UpDownCounter of active requests.
//    * `grpc.server.request.size`: Histogram of request sizes in bytes.
//    * `grpc.server.response.size`: Histogram of response sizes in bytes.
//
// All the metrics have the following labels:
//
//    * `goa.method`: The method name as specified in the Goa design.
//    * `goa.service`: The service name as specified in the Goa design.
//    * `net.peer.name`: The peer name.
//    * `rpc.system`: A stream identifying the remoting system (e.g. `grpc`).
//    * `rpc.service`: Name of RPC service.
//    * `rpc.method`: Name of RPC method.
//    * `rpc.status_code`: The response status code.
//
// Errors collecting or serving metrics are logged to the logger in the context
// if any. The metrics are exp
func UnaryServerInterceptor(ctx context.Context, svc string, opts ...Option) grpc.UnaryServerInterceptor {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	durations := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        MetricRPCDuration,
		Help:        "Histogram of request durations in milliseconds.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
		Buckets:     options.durationBuckets,
	}, RPCLabels)
	options.registerer.MustRegister(durations)

	reqSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        MetricRPCRequestSize,
		Help:        "Histogram of request sizes in bytes.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
		Buckets:     options.requestSizeBuckets,
	}, RPCLabels)
	options.registerer.MustRegister(reqSizes)

	respSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        MetricRPCResponseSize,
		Help:        "Histogram of response sizes in bytes.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
		Buckets:     options.responseSizeBuckets,
	}, RPCLabels)
	options.registerer.MustRegister(respSizes)

	activeReqs := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        MetricRPCActiveRequests,
		Help:        "Gauge of active requests.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
	}, RPCActiveRequestsLabels)
	options.registerer.MustRegister(activeReqs)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		labels := prometheus.Labels{LabelRPCMethod: info.FullMethod}
		if p, ok := peer.FromContext(ctx); ok {
			labels[LabelPeerAddr] = p.Addr.String()
		}
		activeReqs.With(labels).Add(1)
		defer activeReqs.With(labels).Sub(1)

		now := time.Now()
		resp, err := handler(ctx, req)

		labels[LabelGoaMethod] = methodFromCtx(ctx)
		st, _ := status.FromError(err)
		labels[LabelRPCStatusCode] = strconv.Itoa(int(st.Code()))
		durations.With(labels).Observe(float64(timeSince(now)) / float64(time.Millisecond))
		if msg, ok := req.(proto.Message); ok {
			reqSizes.With(labels).Observe(float64(proto.Size(msg)))
		}
		if msg, ok := resp.(proto.Message); ok {
			respSizes.With(labels).Observe(float64(proto.Size(msg)))
		}

		return resp, err
	}
}

// StreamServerInterceptor creates a gRPC stream server interceptor that instruments the
// requests. The returned interceptor adds the following metrics:
//
//    * `grpc.server.active_requests`: UpDownCounter of active requests.
//    * `grpc.server.request.size`: Histogram of request sizes in bytes.
//    * `grpc.server.response.size`: Histogram of response sizes in bytes.
//
// All the metrics have the following labels:
//
//    * `goa.method`: The method name as specified in the Goa design.
//    * `goa.service`: The service name as specified in the Goa design.
//    * `net.peer.name`: The peer name.
//    * `rpc.system`: A stream identifying the remoting system (e.g. `grpc`).
//    * `rpc.service`: Name of RPC service.
//    * `rpc.method`: Name of RPC method.
//    * `rpc.status_code`: The response status code.
//
// Errors collecting or serving metrics are logged to the logger in the context
// if any.
func StreamServerInterceptor(ctx context.Context, svc string, opts ...Option) grpc.StreamServerInterceptor {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	durations := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        MetricRPCDuration,
		Help:        "Histogram of request durations in milliseconds.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
		Buckets:     options.durationBuckets,
	}, RPCLabels)
	options.registerer.MustRegister(durations)

	reqSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        MetricRPCRequestSize,
		Help:        "Histogram of request sizes in bytes.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
		Buckets:     options.requestSizeBuckets,
	}, RPCLabels)
	options.registerer.MustRegister(reqSizes)

	respSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        MetricRPCResponseSize,
		Help:        "Histogram of response sizes in bytes.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
		Buckets:     options.responseSizeBuckets,
	}, RPCLabels)
	options.registerer.MustRegister(respSizes)

	activeReqs := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        MetricRPCActiveRequests,
		Help:        "Gauge of active requests.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
	}, RPCActiveRequestsLabels)
	options.registerer.MustRegister(activeReqs)

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		labels := prometheus.Labels{LabelRPCMethod: info.FullMethod}
		if p, ok := peer.FromContext(ctx); ok {
			labels[LabelPeerAddr] = p.Addr.String()
		}
		activeReqs.With(labels).Add(1)
		defer activeReqs.With(labels).Sub(1)

		now := time.Now()
		wrapper := streamWrapper{stream, labels, reqSizes, respSizes}
		err := handler(srv, &wrapper)

		labels[LabelGoaMethod] = methodFromCtx(ctx)
		st, _ := status.FromError(err)
		labels[LabelRPCStatusCode] = strconv.Itoa(int(st.Code()))

		durations.With(labels).Observe(float64(time.Since(now)) / float64(time.Millisecond))

		return err
	}
}

func (s *streamWrapper) RecvMsg(m interface{}) error {
	if msg, ok := m.(proto.Message); ok {
		s.reqSizes.With(s.labels).Observe(float64(proto.Size(msg)))
	}
	return s.ServerStream.RecvMsg(m)
}

func (s *streamWrapper) SendMsg(m interface{}) error {
	if msg, ok := m.(proto.Message); ok {
		s.respSizes.With(s.labels).Observe(float64(proto.Size(msg)))
	}
	return s.ServerStream.SendMsg(m)
}
