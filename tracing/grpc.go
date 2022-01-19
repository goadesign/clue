package tracing

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
)

// UnaryServerTrace returns an OpenTelemetry UnaryServerInterceptor configured
// to export traces to AWS X-Ray.
func UnaryServerTrace(provider TracerProvider) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		handler = addRequestIDGRPCUnary(handler)
		ui := otelgrpc.UnaryServerInterceptor(
			otelgrpc.WithPropagators(xray.Propagator{}),
			otelgrpc.WithTracerProvider(provider))
		return ui(ctx, req, info, handler)
	}
}

// StreamServerTrace returns an OpenTelemetry StreamServerInterceptor configured
// to export traces to AWS X-Ray.
func StreamServerTrace(provider TracerProvider) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		handler = addRequestIDGRPCStream(handler)
		si := otelgrpc.StreamServerInterceptor(
			otelgrpc.WithPropagators(xray.Propagator{}),
			otelgrpc.WithTracerProvider(provider))
		return si(srv, stream, info, handler)
	}
}

// UnaryClientTrace returns an OpenTelemetry UnaryClientInterceptor configured
// to export traces to AWS X-Ray.
func UnaryClientTrace(provider TracerProvider) grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor(
		otelgrpc.WithPropagators(xray.Propagator{}),
		otelgrpc.WithTracerProvider(provider))
}

// StreamClientTrace returns an OpenTelemetry StreamClientInterceptor configured
// to export traces to AWS X-Ray.
func StreamClientTrace(provider TracerProvider) grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor(
		otelgrpc.WithPropagators(xray.Propagator{}),
		otelgrpc.WithTracerProvider(provider))
}

// addRequestIDGRPCUnary is a middleware that adds the request ID to the current span
// attributes.
func addRequestIDGRPCUnary(h grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		requestID := ctx.Value(middleware.RequestIDKey)
		if requestID == nil {
			return h(ctx, req)
		}
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String(AttributeRequestID, requestID.(string)))
		return h(ctx, req)
	}
}

// addRequestIDGRPCStream is a middleware that adds the request ID to the current span
// attributes.
func addRequestIDGRPCStream(h grpc.StreamHandler) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		requestID := stream.Context().Value(middleware.RequestIDKey)
		if requestID == nil {
			return h(srv, stream)
		}
		span := trace.SpanFromContext(stream.Context())
		span.SetAttributes(attribute.String(AttributeRequestID, requestID.(string)))
		return h(srv, stream)
	}
}
