package instrument

import (
	"context"
	"strconv"
	"strings"
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

// UnaryServerInterceptor creates a gRPC unary server interceptor that instruments the
// requests. The context must have been initialized with instrument.Context. The
// returned interceptor adds the following metrics:
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
// if any.
func UnaryServerInterceptor(ctx context.Context) grpc.UnaryServerInterceptor {
	b := ctx.Value(stateBagKey)
	if b == nil {
		panic("initialize context with Context first")
	}
	metrics := b.(*stateBag).GRPCMetrics()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		service, method := parseGRPCFullMethodName(info.FullMethod)
		labels := prometheus.Labels{labelRPCMethod: method, labelRPCService: service}
		if p, ok := peer.FromContext(ctx); ok {
			ip, port := parseAddr(p.Addr.String())
			labels[labelPeerIP] = ip
			labels[labelPeerPort] = port
		} else {
			labels[labelPeerIP] = ""
			labels[labelPeerPort] = ""
		}
		metrics.ActiveRequests.With(labels).Add(1)
		defer metrics.ActiveRequests.With(labels).Sub(1)

		now := time.Now()
		resp, err := handler(ctx, req)

		st, _ := status.FromError(err)
		labels[labelRPCStatusCode] = strconv.Itoa(int(st.Code()))
		metrics.Durations.With(labels).Observe(float64(timeSince(now)) / float64(time.Millisecond))
		if msg, ok := req.(proto.Message); ok {
			metrics.RequestSizes.With(labels).Observe(float64(proto.Size(msg)))
		}
		if msg, ok := resp.(proto.Message); ok {
			metrics.ResponseSizes.With(labels).Observe(float64(proto.Size(msg)))
		}

		return resp, err
	}
}

// StreamServerInterceptor creates a gRPC stream server interceptor that instruments the
// requests. The context must have been initialized with Context. The returned
// interceptor adds the following metrics:
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
func StreamServerInterceptor(ctx context.Context) grpc.StreamServerInterceptor {
	b := ctx.Value(stateBagKey)
	if b == nil {
		panic("metrics not found in context, initialize context with Context first")
	}
	metrics := b.(*stateBag).GRPCMetrics()

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service, method := parseGRPCFullMethodName(info.FullMethod)
		labels := prometheus.Labels{labelRPCMethod: method, labelRPCService: service}
		if p, ok := peer.FromContext(stream.Context()); ok {
			ip, port := parseAddr(p.Addr.String())
			labels[labelPeerIP] = ip
			labels[labelPeerPort] = port
		} else {
			labels[labelPeerIP] = ""
			labels[labelPeerPort] = ""
		}
		metrics.ActiveRequests.With(labels).Add(1)
		defer metrics.ActiveRequests.With(labels).Sub(1)

		now := time.Now()
		wrapper := streamWrapper{stream, labels, metrics.StreamMessageSizes, metrics.StreamResultSizes}
		err := handler(srv, &wrapper)

		st, _ := status.FromError(err)
		labels[labelRPCStatusCode] = strconv.Itoa(int(st.Code()))

		metrics.Durations.With(labels).Observe(float64(timeSince(now)) / float64(time.Millisecond))

		return err
	}
}

func (s *streamWrapper) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	if msg, ok := m.(proto.Message); ok {
		s.reqSizes.With(s.labels).Observe(float64(proto.Size(msg)))
	}
	return nil
}

func (s *streamWrapper) SendMsg(m interface{}) error {
	if msg, ok := m.(proto.Message); ok {
		s.respSizes.With(s.labels).Observe(float64(proto.Size(msg)))
	}
	return s.ServerStream.SendMsg(m)
}

func parseAddr(addr string) (ip, port string) {
	if addr == "" {
		return "", ""
	}
	if addr[0] == ':' {
		return "", addr[1:]
	}
	if idx := strings.LastIndex(addr, ":"); idx > 0 {
		return addr[:idx], addr[idx+1:]
	}
	return addr, ""
}

func parseGRPCFullMethodName(fullMethodName string) (serviceName, methodName string) {
	if idx := strings.LastIndex(fullMethodName, "."); idx >= 0 {
		fullMethodName = fullMethodName[idx+1:]
	}
	if idx := strings.LastIndex(fullMethodName, "/"); idx > 0 {
		return fullMethodName[:idx], fullMethodName[idx+1:]
	}
	return fullMethodName, ""
}
