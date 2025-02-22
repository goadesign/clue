// Code generated with goa v3.20.0, DO NOT EDIT.
//
// tester protocol buffer definition
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/tester/design -o
// services/tester

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: goagen_tester_tester.proto

package weather_testerpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Tester_TestAll_FullMethodName        = "/weather_tester.Tester/TestAll"
	Tester_TestSmoke_FullMethodName      = "/weather_tester.Tester/TestSmoke"
	Tester_TestForecaster_FullMethodName = "/weather_tester.Tester/TestForecaster"
	Tester_TestLocator_FullMethodName    = "/weather_tester.Tester/TestLocator"
)

// TesterClient is the client API for Tester service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The Weather System Tester Service is used to manage the integration testing
// of the weater system
type TesterClient interface {
	// Runs all tests in the iam system
	TestAll(ctx context.Context, in *TestAllRequest, opts ...grpc.CallOption) (*TestAllResponse, error)
	// Runs smoke tests in the iam system
	TestSmoke(ctx context.Context, in *TestSmokeRequest, opts ...grpc.CallOption) (*TestSmokeResponse, error)
	// Runs tests for the forecaster service
	TestForecaster(ctx context.Context, in *TestForecasterRequest, opts ...grpc.CallOption) (*TestForecasterResponse, error)
	// Runs tests for the locator service
	TestLocator(ctx context.Context, in *TestLocatorRequest, opts ...grpc.CallOption) (*TestLocatorResponse, error)
}

type testerClient struct {
	cc grpc.ClientConnInterface
}

func NewTesterClient(cc grpc.ClientConnInterface) TesterClient {
	return &testerClient{cc}
}

func (c *testerClient) TestAll(ctx context.Context, in *TestAllRequest, opts ...grpc.CallOption) (*TestAllResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestAllResponse)
	err := c.cc.Invoke(ctx, Tester_TestAll_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testerClient) TestSmoke(ctx context.Context, in *TestSmokeRequest, opts ...grpc.CallOption) (*TestSmokeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestSmokeResponse)
	err := c.cc.Invoke(ctx, Tester_TestSmoke_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testerClient) TestForecaster(ctx context.Context, in *TestForecasterRequest, opts ...grpc.CallOption) (*TestForecasterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestForecasterResponse)
	err := c.cc.Invoke(ctx, Tester_TestForecaster_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testerClient) TestLocator(ctx context.Context, in *TestLocatorRequest, opts ...grpc.CallOption) (*TestLocatorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestLocatorResponse)
	err := c.cc.Invoke(ctx, Tester_TestLocator_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TesterServer is the server API for Tester service.
// All implementations must embed UnimplementedTesterServer
// for forward compatibility.
//
// The Weather System Tester Service is used to manage the integration testing
// of the weater system
type TesterServer interface {
	// Runs all tests in the iam system
	TestAll(context.Context, *TestAllRequest) (*TestAllResponse, error)
	// Runs smoke tests in the iam system
	TestSmoke(context.Context, *TestSmokeRequest) (*TestSmokeResponse, error)
	// Runs tests for the forecaster service
	TestForecaster(context.Context, *TestForecasterRequest) (*TestForecasterResponse, error)
	// Runs tests for the locator service
	TestLocator(context.Context, *TestLocatorRequest) (*TestLocatorResponse, error)
	mustEmbedUnimplementedTesterServer()
}

// UnimplementedTesterServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTesterServer struct{}

func (UnimplementedTesterServer) TestAll(context.Context, *TestAllRequest) (*TestAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestAll not implemented")
}
func (UnimplementedTesterServer) TestSmoke(context.Context, *TestSmokeRequest) (*TestSmokeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestSmoke not implemented")
}
func (UnimplementedTesterServer) TestForecaster(context.Context, *TestForecasterRequest) (*TestForecasterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestForecaster not implemented")
}
func (UnimplementedTesterServer) TestLocator(context.Context, *TestLocatorRequest) (*TestLocatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestLocator not implemented")
}
func (UnimplementedTesterServer) mustEmbedUnimplementedTesterServer() {}
func (UnimplementedTesterServer) testEmbeddedByValue()                {}

// UnsafeTesterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TesterServer will
// result in compilation errors.
type UnsafeTesterServer interface {
	mustEmbedUnimplementedTesterServer()
}

func RegisterTesterServer(s grpc.ServiceRegistrar, srv TesterServer) {
	// If the following call pancis, it indicates UnimplementedTesterServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Tester_ServiceDesc, srv)
}

func _Tester_TestAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TesterServer).TestAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Tester_TestAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TesterServer).TestAll(ctx, req.(*TestAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tester_TestSmoke_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestSmokeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TesterServer).TestSmoke(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Tester_TestSmoke_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TesterServer).TestSmoke(ctx, req.(*TestSmokeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tester_TestForecaster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestForecasterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TesterServer).TestForecaster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Tester_TestForecaster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TesterServer).TestForecaster(ctx, req.(*TestForecasterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tester_TestLocator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestLocatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TesterServer).TestLocator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Tester_TestLocator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TesterServer).TestLocator(ctx, req.(*TestLocatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Tester_ServiceDesc is the grpc.ServiceDesc for Tester service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tester_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "weather_tester.Tester",
	HandlerType: (*TesterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TestAll",
			Handler:    _Tester_TestAll_Handler,
		},
		{
			MethodName: "TestSmoke",
			Handler:    _Tester_TestSmoke_Handler,
		},
		{
			MethodName: "TestForecaster",
			Handler:    _Tester_TestForecaster_Handler,
		},
		{
			MethodName: "TestLocator",
			Handler:    _Tester_TestLocator_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "goagen_tester_tester.proto",
}
