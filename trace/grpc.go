package trace

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns an OpenTelemetry UnaryServerInterceptor. It
// panics if the context has not been initialized with Context.
func UnaryServerInterceptor(traceCtx context.Context) grpc.UnaryServerInterceptor {
	state := traceCtx.Value(stateKey)
	if state == nil {
		panic(errContextMissing)
	}
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		handler = initTracingContextGRPCUnary(traceCtx, handler)
		handler = addRequestIDGRPCUnary(handler)
		ui := otelgrpc.UnaryServerInterceptor(
			otelgrpc.WithTracerProvider(state.(*stateBag).provider),
			otelgrpc.WithPropagators(state.(*stateBag).propagator))
		return ui(ctx, req, info, handler)
	}
}

// StreamServerInterceptor returns an OpenTelemetry StreamServerInterceptor. It
// panics if the context has not been initialized with Context.
func StreamServerInterceptor(traceCtx context.Context) grpc.StreamServerInterceptor {
	state := traceCtx.Value(stateKey)
	if state == nil {
		panic(errContextMissing)
	}
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		handler = initTracingContextGRPCStream(traceCtx, handler)
		handler = addRequestIDGRPCStream(handler)
		si := otelgrpc.StreamServerInterceptor(
			otelgrpc.WithTracerProvider(state.(*stateBag).provider),
			otelgrpc.WithPropagators(state.(*stateBag).propagator),
		)
		return si(srv, stream, info, handler)
	}
}

// UnaryClientInterceptor returns an OpenTelemetry UnaryClientInterceptor. It
// panics if the context has not been initialized with Context.
func UnaryClientInterceptor(traceCtx context.Context) grpc.UnaryClientInterceptor {
	state := traceCtx.Value(stateKey)
	if state == nil {
		panic(errContextMissing)
	}
	return otelgrpc.UnaryClientInterceptor(
		otelgrpc.WithTracerProvider(state.(*stateBag).provider),
		otelgrpc.WithPropagators(state.(*stateBag).propagator))
}

// StreamClientInterceptor returns an OpenTelemetry StreamClientInterceptor. It
// panics if the context has not been initialized with Context.
func StreamClientInterceptor(traceCtx context.Context) grpc.StreamClientInterceptor {
	state := traceCtx.Value(stateKey)
	if state == nil {
		panic(errContextMissing)
	}
	return otelgrpc.StreamClientInterceptor(
		otelgrpc.WithTracerProvider(state.(*stateBag).provider),
		otelgrpc.WithPropagators(state.(*stateBag).propagator))
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

// initTracingContextGRPCUnary is a unary interceptor that initializes the
// tracing context.
func initTracingContextGRPCUnary(traceCtx context.Context, h grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if IsTraced(ctx) {
			ctx = withTracing(traceCtx, ctx)
		}
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

// initTracingContextGRPCStream is a stream interceptor that initializes the
// tracing context.
func initTracingContextGRPCStream(traceCtx context.Context, h grpc.StreamHandler) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		if IsTraced(stream.Context()) {
			ctx := withTracing(traceCtx, stream.Context())
			stream = &streamWithContext{ctx: ctx, ServerStream: stream}
		}
		return h(srv, stream)
	}
}

type streamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *streamWithContext) Context() context.Context {
	return s.ctx
}
