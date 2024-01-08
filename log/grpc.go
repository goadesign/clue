package log

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// GRPCLogOption is a function that applies a configuration option
	// to a GRPC interceptor logger.
	GRPCLogOption func(*grpcOptions)

	// GRPCClientLogOption is a function that applies a configuration option
	// to a GRPC client interceptor logger.
	// Deprecated: Use GRPCLogOption instead.
	GRPCClientLogOption = GRPCLogOption

	grpcOptions struct {
		iserr              func(codes.Code) bool
		disableCallLogging bool
		disableCallID      bool
	}
)

// Be nice to tests
var shortID = randShortID

// UnaryServerInterceptor returns a unary interceptor that performs two tasks:
// 1. Enriches the request context with the logger specified in logCtx.
// 2. Logs details of the unary call, unless the WithDisableCallLogging option is provided.
// UnaryServerInterceptor panics if logCtx was not created with Context.
func UnaryServerInterceptor(logCtx context.Context, opts ...GRPCLogOption) grpc.UnaryServerInterceptor {
	MustContainLogger(logCtx)
	o := defaultGRPCOptions()
	for _, opt := range opts {
		opt(o)
	}
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = WithContext(ctx, logCtx)
		if !o.disableCallID {
			ctx = With(ctx, KV{RequestIDKey, shortID()})
		}
		if o.disableCallLogging {
			return handler(ctx, req)
		}
		then := time.Now()
		svcKV := KV{K: GRPCServiceKey, V: path.Dir(info.FullMethod)[1:]}
		methKV := KV{K: GRPCMethodKey, V: path.Base(info.FullMethod)}
		Print(ctx, KV{MessageKey, "start"}, svcKV, methKV)

		res, err := handler(ctx, req)

		stat, _ := status.FromError(err)
		ms := timeSince(then).Milliseconds()
		codeKV := KV{K: GRPCCodeKey, V: stat.Code()}
		durKV := KV{K: GRPCDurationKey, V: ms}
		if o.iserr(stat.Code()) {
			statKV := KV{K: GRPCStatusKey, V: stat.Message()}
			Error(ctx, err, svcKV, methKV, statKV, codeKV, durKV)
			return res, err
		}
		Print(ctx, KV{MessageKey, "end"}, svcKV, methKV, codeKV, durKV)
		return res, err
	}
}

// StreamServerInterceptor returns a stream interceptor that performs two tasks:
// 1. Enriches the request context with the logger specified in logCtx.
// 2. Logs details of the stream call, unless the WithDisableCallLogging option is provided.
// StreamServerInterceptor panics if logCtx was not created with Context.
func StreamServerInterceptor(logCtx context.Context, opts ...GRPCLogOption) grpc.StreamServerInterceptor {
	MustContainLogger(logCtx)
	o := defaultGRPCOptions()
	for _, opt := range opts {
		opt(o)
	}
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := WithContext(stream.Context(), logCtx)
		if !o.disableCallID {
			ctx = With(ctx, KV{RequestIDKey, shortID()})
		}
		stream = &streamWithContext{stream, ctx}
		if o.disableCallLogging {
			return handler(srv, stream)
		}
		then := time.Now()
		svcKV := KV{K: GRPCServiceKey, V: path.Dir(info.FullMethod)[1:]}
		methKV := KV{K: GRPCMethodKey, V: path.Base(info.FullMethod)}
		Print(ctx, KV{MessageKey, "start"}, svcKV, methKV)

		err := handler(srv, stream)

		stat, _ := status.FromError(err)
		ms := timeSince(then).Milliseconds()
		codeKV := KV{K: GRPCCodeKey, V: stat.Code()}
		durKV := KV{K: GRPCDurationKey, V: ms}
		if o.iserr(stat.Code()) {
			statKV := KV{K: GRPCStatusKey, V: stat.Message()}
			Error(ctx, err, svcKV, methKV, statKV, codeKV, durKV)
			return err
		}
		Print(ctx, KV{MessageKey, "end"}, svcKV, methKV, codeKV, durKV)
		return err
	}
}

// UnaryClientInterceptor returns a unary interceptor that logs the request with
// the logger contained in the request context if any.
func UnaryClientInterceptor(opts ...GRPCLogOption) grpc.UnaryClientInterceptor {
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
		svcKV := KV{K: GRPCServiceKey, V: path.Dir(fullmethod)[1:]}
		methKV := KV{K: GRPCMethodKey, V: path.Base(fullmethod)}
		Print(ctx, KV{K: MessageKey, V: "start"}, svcKV, methKV)

		err := invoker(ctx, fullmethod, req, reply, cc, opts...)

		stat, _ := status.FromError(err)
		ms := timeSince(then).Milliseconds()
		codeKV := KV{K: GRPCCodeKey, V: stat.Code()}
		durKV := KV{K: GRPCDurationKey, V: ms}
		if o.iserr(stat.Code()) {
			statKV := KV{K: GRPCStatusKey, V: stat.Message()}
			Error(ctx, err, svcKV, methKV, statKV, codeKV, durKV)
			return err
		}
		Print(ctx, KV{K: MessageKey, V: "end"}, svcKV, methKV, codeKV, durKV)
		return err
	}
}

// StreamClientInterceptor returns a stream interceptor that logs the request
// with the logger contained in the request context if any.
func StreamClientInterceptor(opts ...GRPCLogOption) grpc.StreamClientInterceptor {
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
		svcKV := KV{K: GRPCServiceKey, V: path.Dir(fullmethod)[1:]}
		methKV := KV{K: GRPCMethodKey, V: path.Base(fullmethod)}
		Print(ctx, KV{K: MessageKey, V: "start"}, svcKV, methKV)

		stream, err := streamer(ctx, desc, cc, fullmethod, opts...)

		stat, _ := status.FromError(err)
		ms := timeSince(then).Milliseconds()
		codeKV := KV{K: GRPCCodeKey, V: stat.Code()}
		durKV := KV{K: GRPCDurationKey, V: ms}
		if o.iserr(stat.Code()) {
			statKV := KV{K: GRPCStatusKey, V: stat.Message()}
			Error(ctx, err, svcKV, methKV, statKV, codeKV, durKV)
			return stream, err
		}
		Print(ctx, KV{K: MessageKey, V: "end"}, svcKV, methKV, codeKV, durKV)
		return stream, err
	}
}

// WithErrorFunc returns a GRPC logger option that configures the logger to
// consider the given function to determine if a GRPC status code is an error.
func WithErrorFunc(iserr func(codes.Code) bool) GRPCLogOption {
	return func(o *grpcOptions) {
		o.iserr = iserr
	}
}

// WithDisableCallLogging returns a GRPC logger option that disables call
// logging.
func WithDisableCallLogging() GRPCLogOption {
	return func(o *grpcOptions) {
		o.disableCallLogging = true
	}
}

// WithDisableCallID returns a GRPC logger option that disables the
// generation of request IDs.
func WithDisableCallID() GRPCLogOption {
	return func(o *grpcOptions) {
		o.disableCallID = true
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

// randShortID produces a "unique" 6 bytes long string.
// This algorithm favors simplicity and efficiency over true uniqueness.
func randShortID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b) // nolint: errcheck
	return base64.RawURLEncoding.EncodeToString(b)
}
