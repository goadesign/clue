// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.3
// source: test.proto

package testpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TestClient is the client API for Test service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestClient interface {
	// GrpcMethod implements grpc_method.
	GrpcMethod(ctx context.Context, in *GrpcMethodRequest, opts ...grpc.CallOption) (*GrpcMethodResponse, error)
	// GrpcStream implements grpc_stream.
	GrpcStream(ctx context.Context, opts ...grpc.CallOption) (Test_GrpcStreamClient, error)
}

type testClient struct {
	cc grpc.ClientConnInterface
}

func NewTestClient(cc grpc.ClientConnInterface) TestClient {
	return &testClient{cc}
}

func (c *testClient) GrpcMethod(ctx context.Context, in *GrpcMethodRequest, opts ...grpc.CallOption) (*GrpcMethodResponse, error) {
	out := new(GrpcMethodResponse)
	err := c.cc.Invoke(ctx, "/test.Test/GrpcMethod", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testClient) GrpcStream(ctx context.Context, opts ...grpc.CallOption) (Test_GrpcStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Test_ServiceDesc.Streams[0], "/test.Test/GrpcStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &testGrpcStreamClient{stream}
	return x, nil
}

type Test_GrpcStreamClient interface {
	Send(*GrpcStreamStreamingRequest) error
	Recv() (*GrpcStreamResponse, error)
	grpc.ClientStream
}

type testGrpcStreamClient struct {
	grpc.ClientStream
}

func (x *testGrpcStreamClient) Send(m *GrpcStreamStreamingRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testGrpcStreamClient) Recv() (*GrpcStreamResponse, error) {
	m := new(GrpcStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestServer is the server API for Test service.
// All implementations must embed UnimplementedTestServer
// for forward compatibility
type TestServer interface {
	// GrpcMethod implements grpc_method.
	GrpcMethod(context.Context, *GrpcMethodRequest) (*GrpcMethodResponse, error)
	// GrpcStream implements grpc_stream.
	GrpcStream(Test_GrpcStreamServer) error
	mustEmbedUnimplementedTestServer()
}

// UnimplementedTestServer must be embedded to have forward compatible implementations.
type UnimplementedTestServer struct {
}

func (UnimplementedTestServer) GrpcMethod(context.Context, *GrpcMethodRequest) (*GrpcMethodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GrpcMethod not implemented")
}
func (UnimplementedTestServer) GrpcStream(Test_GrpcStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GrpcStream not implemented")
}
func (UnimplementedTestServer) mustEmbedUnimplementedTestServer() {}

// UnsafeTestServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestServer will
// result in compilation errors.
type UnsafeTestServer interface {
	mustEmbedUnimplementedTestServer()
}

func RegisterTestServer(s grpc.ServiceRegistrar, srv TestServer) {
	s.RegisterService(&Test_ServiceDesc, srv)
}

func _Test_GrpcMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcMethodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServer).GrpcMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.Test/GrpcMethod",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServer).GrpcMethod(ctx, req.(*GrpcMethodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Test_GrpcStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TestServer).GrpcStream(&testGrpcStreamServer{stream})
}

type Test_GrpcStreamServer interface {
	Send(*GrpcStreamResponse) error
	Recv() (*GrpcStreamStreamingRequest, error)
	grpc.ServerStream
}

type testGrpcStreamServer struct {
	grpc.ServerStream
}

func (x *testGrpcStreamServer) Send(m *GrpcStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *testGrpcStreamServer) Recv() (*GrpcStreamStreamingRequest, error) {
	m := new(GrpcStreamStreamingRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Test_ServiceDesc is the grpc.ServiceDesc for Test service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Test_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "test.Test",
	HandlerType: (*TestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GrpcMethod",
			Handler:    _Test_GrpcMethod_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GrpcStream",
			Handler:       _Test_GrpcStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "test.proto",
}
