package trace

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns an OpenTelemetry UnaryServerInterceptor configured
// to export traces to AWS X-Ray. It panics if the context has not been
// initialized with Context.
func UnaryServerInterceptor(ctx context.Context) grpc.UnaryServerInterceptor {
	s := ctx.Value(stateKey)
	if s == nil {
		panic(errContextMissing)
	}
	return func(
		reqctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		handler = initTracingContextGRPCUnary(ctx, handler)
		handler = addRequestIDGRPCUnary(handler)
		ui := otelgrpc.UnaryServerInterceptor(
			otelgrpc.WithTracerProvider(s.(*stateBag).provider))
		return ui(reqctx, req, info, handler)
	}
}

// StreamServerInterceptor returns an OpenTelemetry StreamServerInterceptor configured
// to export traces to AWS X-Ray. It panics if the context has not been
// initialized with Context.
func StreamServerInterceptor(ctx context.Context) grpc.StreamServerInterceptor {
	s := ctx.Value(stateKey)
	if s == nil {
		panic(errContextMissing)
	}
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		handler = initTracingContextGRPCStream(ctx, handler)
		handler = addRequestIDGRPCStream(handler)
		si := otelgrpc.StreamServerInterceptor(
			otelgrpc.WithTracerProvider(s.(*stateBag).provider))
		return si(srv, stream, info, handler)
	}
}

// UnaryClientInterceptor returns an OpenTelemetry UnaryClientInterceptor configured
// to export traces to AWS X-Ray. It panics if the context has not been
// initialized with Context.
func UnaryClientInterceptor(ctx context.Context) grpc.UnaryClientInterceptor {
	s := ctx.Value(stateKey)
	if s == nil {
		panic(errContextMissing)
	}
	return otelgrpc.UnaryClientInterceptor(
		otelgrpc.WithTracerProvider(s.(*stateBag).provider))
}

// StreamClientInterceptor returns an OpenTelemetry StreamClientInterceptor configured
// to export traces to AWS X-Ray. It panics if the context has not been
// initialized with Context.
func StreamClientInterceptor(ctx context.Context) grpc.StreamClientInterceptor {
	s := ctx.Value(stateKey)
	if s == nil {
		panic(errContextMissing)
	}
	return otelgrpc.StreamClientInterceptor(
		otelgrpc.WithTracerProvider(s.(*stateBag).provider))
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
		s := traceCtx.Value(stateKey).(*stateBag)
		ctx = withProvider(ctx, s.provider)
		setActiveSpans(ctx, []trace.Span{trace.SpanFromContext(ctx)})
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
		s := traceCtx.Value(stateKey).(*stateBag)
		ctx := withProvider(stream.Context(), s.provider)
		setActiveSpans(ctx, []trace.Span{trace.SpanFromContext(ctx)})
		return h(srv, &streamWithContext{stream, ctx})
	}
}

type streamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *streamWithContext) Context() context.Context {
	return s.ctx
}
