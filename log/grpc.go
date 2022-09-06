package log

import (
	"context"
	"path"
	"time"

	goamiddleware "goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// GRPCClientLogOption is a function that applies a configuration option
	// to a GRPC client interceptor logger.
	GRPCClientLogOption func(*grpcOptions)

	grpcOptions struct {
		iserr func(codes.Code) bool
	}
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
		if reqID := ctx.Value(goamiddleware.RequestIDKey); reqID != nil {
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
		if reqID := ctx.Value(goamiddleware.RequestIDKey); reqID != nil {
			ctx = With(ctx, KV{RequestIDKey, reqID})
		}
		return handler(srv, &streamWithContext{stream, ctx})
	}
}

// UnaryClientInterceptor returns a unary interceptor that logs the request with
// the logger contained in the request context if any.
func UnaryClientInterceptor(opts ...GRPCClientLogOption) grpc.UnaryClientInterceptor {
	o := defaultGRPCOptions()
	for _, opt := range opts {
		opt(o)
	}
	return func(
		ctx context.Context,
		fullmethod string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		then := time.Now()
		service := path.Dir(fullmethod)[1:]
		method := path.Base(fullmethod)
		err := invoker(ctx, fullmethod, req, reply, cc, opts...)
		stat, _ := status.FromError(err)
		ms := timeSince(then).Milliseconds()
		msgKV := KV{K: MessageKey, V: "finished client unary call"}
		svcKV := KV{K: GRPCServiceKey, V: service}
		methKV := KV{K: GRPCMethodKey, V: method}
		codeKV := KV{K: GRPCCodeKey, V: stat.Code()}
		durKV := KV{K: GRPCDurationKey, V: ms}
		if o.iserr(stat.Code()) {
			statKV := KV{K: GRPCStatusKey, V: stat.Message()}
			Error(ctx, err, msgKV, svcKV, methKV, statKV, codeKV, durKV)
			return err
		}
		Print(ctx, msgKV, svcKV, methKV, codeKV, durKV)
		return err
	}
}

// StreamClientInterceptor returns a stream interceptor that logs the request
// with the logger contained in the request context if any.
func StreamClientInterceptor(opts ...GRPCClientLogOption) grpc.StreamClientInterceptor {
	o := defaultGRPCOptions()
	for _, opt := range opts {
		opt(o)
	}
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		fullmethod string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		then := time.Now()
		service := path.Dir(fullmethod)[1:]
		method := path.Base(fullmethod)
		stream, err := streamer(ctx, desc, cc, fullmethod, opts...)
		stat, _ := status.FromError(err)
		ms := timeSince(then).Milliseconds()
		msgKV := KV{K: MessageKey, V: "finished client streaming call"}
		svcKV := KV{K: GRPCServiceKey, V: service}
		methKV := KV{K: GRPCMethodKey, V: method}
		codeKV := KV{K: GRPCCodeKey, V: stat.Code()}
		durKV := KV{K: GRPCDurationKey, V: ms}
		if o.iserr(stat.Code()) {
			statKV := KV{K: GRPCStatusKey, V: stat.Message()}
			Error(ctx, err, msgKV, svcKV, methKV, statKV, codeKV, durKV)
			return stream, err
		}
		Print(ctx, msgKV, svcKV, methKV, codeKV, durKV)
		return stream, err
	}
}

func WithErrorFunc(iserr func(codes.Code) bool) GRPCClientLogOption {
	return func(o *grpcOptions) {
		o.iserr = iserr
	}
}

func defaultGRPCOptions() *grpcOptions {
	return &grpcOptions{
		iserr: func(c codes.Code) bool {
			return c != codes.OK
		},
	}
}

type streamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *streamWithContext) Context() context.Context {
	return s.ctx
}
