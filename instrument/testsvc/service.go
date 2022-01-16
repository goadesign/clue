package testsvc

import (
	"context"

	"github.com/crossnokaye/micro/instrument/testsvc/gen/test"
)

type (
	UnaryFunc      func(context.Context, *Fields) (res *Fields, err error)
	GRPCStreamFunc func(context.Context, test.GrpcStreamingServerStream) (err error)

	// Shadow generated type to avoid dependency creep.
	Fields       test.Fields
	ServerStream test.GrpcStreamingServerStream

	svc struct {
		httpfn   UnaryFunc
		grpcfn   UnaryFunc
		streamfn GRPCStreamFunc
	}
)

func (s *svc) HTTPMethod(ctx context.Context, req *test.Fields) (res *test.Fields, err error) {
	if s.httpfn == nil {
		return
	}

	var r *Fields
	if req != nil {
		r = &Fields{}
		*r = Fields(*req)
	}
	var resp *Fields
	resp, err = s.httpfn(ctx, r)
	if resp != nil {
		res = &test.Fields{}
		*res = test.Fields(*resp)
	}
	return
}

func (s *svc) GrpcMethod(ctx context.Context, req *test.Fields) (res *test.Fields, err error) {
	if s.grpcfn == nil {
		return
	}

	var r *Fields
	if req != nil {
		r = &Fields{}
		*r = Fields(*req)
	}
	var resp *Fields
	resp, err = s.grpcfn(ctx, r)
	if resp != nil {
		res = &test.Fields{}
		*res = test.Fields(*resp)
	}
	return
}

func (s *svc) GrpcStreaming(ctx context.Context, stream test.GrpcStreamingServerStream) (err error) {
	if s.streamfn != nil {
		return s.streamfn(ctx, stream)
	}
	return nil
}
