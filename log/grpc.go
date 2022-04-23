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
	MustContainLogger(logCtx)
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = WithContext(ctx, logCtx)
		if reqID := ctx.Value(middleware.RequestIDKey); reqID != nil {
			ctx = With(ctx, KV{RequestIDKey, reqID})
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a stream interceptor that configures the
// request context with the logger contained in logCtx.  It panics if logCtx
// was not initialized with Context.
func StreamServerInterceptor(logCtx context.Context) grpc.StreamServerInterceptor {
	MustContainLogger(logCtx)
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := WithContext(stream.Context(), logCtx)
		if reqID := ctx.Value(middleware.RequestIDKey); reqID != nil {
			ctx = With(ctx, KV{RequestIDKey, reqID})
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
