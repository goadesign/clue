// Code generated by goa v3.8.3, DO NOT EDIT.
//
// test gRPC server
//
// Command:
// $ goa gen goa.design/clue/internal/testsvc/design

package server

import (
	"context"

	testpb "goa.design/clue/internal/testsvc/gen/grpc/test/pb"
	test "goa.design/clue/internal/testsvc/gen/test"
	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
)

// Server implements the testpb.TestServer interface.
type Server struct {
	GrpcMethodH goagrpc.UnaryHandler
	GrpcStreamH goagrpc.StreamHandler
	testpb.UnimplementedTestServer
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the expr.
type ErrorNamer interface {
	ErrorName() string
}

// GrpcStreamServerStream implements the test.GrpcStreamServerStream interface.
type GrpcStreamServerStream struct {
	stream testpb.Test_GrpcStreamServer
}

// New instantiates the server struct with the test service endpoints.
func New(e *test.Endpoints, uh goagrpc.UnaryHandler, sh goagrpc.StreamHandler) *Server {
	return &Server{
		GrpcMethodH: NewGrpcMethodHandler(e.GrpcMethod, uh),
		GrpcStreamH: NewGrpcStreamHandler(e.GrpcStream, sh),
	}
}

// NewGrpcMethodHandler creates a gRPC handler which serves the "test" service
// "grpc_method" endpoint.
func NewGrpcMethodHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeGrpcMethodRequest, EncodeGrpcMethodResponse)
	}
	return h
}

// GrpcMethod implements the "GrpcMethod" method in testpb.TestServer interface.
func (s *Server) GrpcMethod(ctx context.Context, message *testpb.GrpcMethodRequest) (*testpb.GrpcMethodResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "grpc_method")
	ctx = context.WithValue(ctx, goa.ServiceKey, "test")
	resp, err := s.GrpcMethodH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*testpb.GrpcMethodResponse), nil
}

// NewGrpcStreamHandler creates a gRPC handler which serves the "test" service
// "grpc_stream" endpoint.
func NewGrpcStreamHandler(endpoint goa.Endpoint, h goagrpc.StreamHandler) goagrpc.StreamHandler {
	if h == nil {
		h = goagrpc.NewStreamHandler(endpoint, nil)
	}
	return h
}

// GrpcStream implements the "GrpcStream" method in testpb.TestServer interface.
func (s *Server) GrpcStream(stream testpb.Test_GrpcStreamServer) error {
	ctx := stream.Context()
	ctx = context.WithValue(ctx, goa.MethodKey, "grpc_stream")
	ctx = context.WithValue(ctx, goa.ServiceKey, "test")
	_, err := s.GrpcStreamH.Decode(ctx, nil)
	if err != nil {
		return goagrpc.EncodeError(err)
	}
	ep := &test.GrpcStreamEndpointInput{
		Stream: &GrpcStreamServerStream{stream: stream},
	}
	err = s.GrpcStreamH.Handle(ctx, ep)
	if err != nil {
		return goagrpc.EncodeError(err)
	}
	return nil
}

// Send streams instances of "testpb.GrpcStreamResponse" to the "grpc_stream"
// endpoint gRPC stream.
func (s *GrpcStreamServerStream) Send(res *test.Fields) error {
	v := NewProtoGrpcStreamResponse(res)
	return s.stream.Send(v)
}

// Recv reads instances of "testpb.GrpcStreamStreamingRequest" from the
// "grpc_stream" endpoint gRPC stream.
func (s *GrpcStreamServerStream) Recv() (*test.Fields, error) {
	var res *test.Fields
	v, err := s.stream.Recv()
	if err != nil {
		return res, err
	}
	return NewFields(v), nil
}

func (s *GrpcStreamServerStream) Close() error {
	// nothing to do here
	return nil
}
