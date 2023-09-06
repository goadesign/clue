package debug

import (
	"context"

	"google.golang.org/grpc"

	"goa.design/clue/log"
)

// UnaryServerInterceptor return an interceptor that manages whether debug log
// entries are written. This interceptor should be used in conjunction with the
// MountDebugLogEnabler function.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if debugLogs {
			ctx = log.Context(ctx, log.WithDebug())
		} else {
			ctx = log.Context(ctx, log.WithNoDebug())
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a stream interceptor that manages whether
// debug log entries are written. Note: a change in the debug setting is
// effective only for the next stream request. This interceptor should be used
// in conjunction with the MountDebugLogEnabler function.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		_ *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := stream.Context()
		if debugLogs {
			ctx = log.Context(ctx, log.WithDebug())
		} else {
			ctx = log.Context(ctx, log.WithNoDebug())
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
