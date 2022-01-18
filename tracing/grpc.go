package tracing

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

// UnaryServerTrace returns an OpenTelemetry UnaryServerInterceptor configured
// to export traces to AWS X-Ray.
func UnaryServerTrace(provider *sdktrace.TracerProvider) grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithPropagators(xray.Propagator{}),
		otelgrpc.WithTracerProvider(provider))
}

// StreamServerTrace returns an OpenTelemetry StreamServerInterceptor configured
// to export traces to AWS X-Ray.
func StreamServerTrace(provider *sdktrace.TracerProvider) grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor(
		otelgrpc.WithPropagators(xray.Propagator{}),
		otelgrpc.WithTracerProvider(provider))
}

// UnaryClientTrace returns an OpenTelemetry UnaryClientInterceptor configured
// to export traces to AWS X-Ray.
func UnaryClientTrace(provider *sdktrace.TracerProvider) grpc.UnaryClientInterceptor {
	return otelgrpc.UnaryClientInterceptor(
		otelgrpc.WithPropagators(xray.Propagator{}),
		otelgrpc.WithTracerProvider(provider))
}

// StreamClientTrace returns an OpenTelemetry StreamClientInterceptor configured
// to export traces to AWS X-Ray.
func StreamClientTrace(provider *sdktrace.TracerProvider) grpc.StreamClientInterceptor {
	return otelgrpc.StreamClientInterceptor(
		otelgrpc.WithPropagators(xray.Propagator{}),
		otelgrpc.WithTracerProvider(provider))
}
