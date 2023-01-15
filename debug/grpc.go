package debug

import (
	"context"

	"google.golang.org/grpc"

	"goa.design/clue/log"
)

// UnaryServerInterceptor return an interceptor that manages whether debug log
// entries are written.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if wantDebugEnabled && !debugEnabled {
			ctx = log.Context(ctx, log.WithDebug())
			log.Debug(ctx, log.KV{K: "debug-logs", V: true})
			debugEnabled = true
		} else if !wantDebugEnabled && debugEnabled {
			log.Debug(ctx, log.KV{K: "debug-logs", V: false})
			ctx = log.Context(ctx, log.WithNoDebug())
			debugEnabled = false
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a stream interceptor that manages whether
// debug log entries are written. Note: a change in the debug setting is
// effective only for the next stream request.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := stream.Context()
		if wantDebugEnabled && !debugEnabled {
			ctx = log.Context(ctx, log.WithDebug())
			log.Debug(ctx, log.KV{K: "debug", V: true})
			debugEnabled = true
		} else if !wantDebugEnabled && debugEnabled {
			log.Debug(ctx, log.KV{K: "debug", V: false})
			ctx = log.Context(ctx, log.WithNoDebug())
			debugEnabled = false
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
