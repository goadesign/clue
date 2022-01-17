package testsvc

import (
	"context"

	"github.com/crossnokaye/micro/instrument/testsvc/gen/test"
)

type (
	UnaryFunc  func(context.Context, *Fields) (res *Fields, err error)
	StreamFunc func(context.Context, Stream) (err error)

	// Shadow generated type to avoid dependency creep.
	Fields test.Fields
	Stream interface {
		Send(*Fields) error
		Recv() (*Fields, error)
		Close() error
	}

	svc struct {
		httpfn   UnaryFunc
		grpcfn   UnaryFunc
		streamfn StreamFunc
	}

	adapter struct {
		stream test.GrpcStreamServerStream
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

func (s *svc) GrpcStream(ctx context.Context, stream test.GrpcStreamServerStream) (err error) {
	if s.streamfn != nil {
		return s.streamfn(ctx, adapter{stream})
	}
	return nil
}

func (a adapter) Send(fields *Fields) error {
	var f *test.Fields
	if fields != nil {
		f = &test.Fields{}
		*f = test.Fields(*fields)
	}
	return a.stream.Send(f)
}

func (a adapter) Recv() (*Fields, error) {
	f, err := a.stream.Recv()
	if f == nil {
		return nil, err
	}
	fields := &Fields{}
	*fields = Fields(*f)
	return fields, err
}

func (a adapter) Close() error {
	return a.stream.Close()
}
