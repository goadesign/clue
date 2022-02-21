package log

import (
	"context"

	"goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor return an interceptor that configured the request
// context with the logger contained in logCtx.  It panics if logCtx was not
// initialized with Context.
func UnaryServerInterceptor(logCtx context.Context) grpc.UnaryServerInterceptor {
	l := logCtx.Value(ctxLogger)
	if l == nil {
		panic("log.Init called without log.Context")
	}
	logger := l.(*logger)
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = context.WithValue(ctx, ctxLogger, logger)
		if reqID := ctx.Value(middleware.RequestIDKey); reqID != nil {
			ctx = With(ctx, KV{"requestID", reqID})
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a stream interceptor that configures the
// request context with the logger contained in logCtx.  It panics if logCtx
// was not initialized with Context.
func StreamServerInterceptor(logCtx context.Context) grpc.StreamServerInterceptor {
	l := logCtx.Value(ctxLogger)
	if l == nil {
		panic("log.Init called without log.Context")
	}
	logger := l.(*logger)
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := context.WithValue(stream.Context(), ctxLogger, logger)
		if reqID := ctx.Value(middleware.RequestIDKey); reqID != nil {
			ctx = With(ctx, KV{"requestID", reqID})
		}
		return handler(srv, &streamWithContext{stream, ctx})
	}
}

type streamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *streamWithContext) Context() context.Context {
	return s.ctx
}
