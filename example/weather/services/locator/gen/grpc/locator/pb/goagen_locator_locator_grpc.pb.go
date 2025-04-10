// Code generated with goa v3.20.0, DO NOT EDIT.
//
// locator protocol buffer definition
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/locator/design -o
// services/locator

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: goagen_locator_locator.proto

package locatorpb

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
	Locator_GetLocation_FullMethodName = "/locator.Locator/GetLocation"
)

// LocatorClient is the client API for Locator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Public HTTP frontend
type LocatorClient interface {
	// Retrieve location information for a given IP address
	GetLocation(ctx context.Context, in *GetLocationRequest, opts ...grpc.CallOption) (*GetLocationResponse, error)
}

type locatorClient struct {
	cc grpc.ClientConnInterface
}

func NewLocatorClient(cc grpc.ClientConnInterface) LocatorClient {
	return &locatorClient{cc}
}

func (c *locatorClient) GetLocation(ctx context.Context, in *GetLocationRequest, opts ...grpc.CallOption) (*GetLocationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetLocationResponse)
	err := c.cc.Invoke(ctx, Locator_GetLocation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LocatorServer is the server API for Locator service.
// All implementations must embed UnimplementedLocatorServer
// for forward compatibility.
//
// Public HTTP frontend
type LocatorServer interface {
	// Retrieve location information for a given IP address
	GetLocation(context.Context, *GetLocationRequest) (*GetLocationResponse, error)
	mustEmbedUnimplementedLocatorServer()
}

// UnimplementedLocatorServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLocatorServer struct{}

func (UnimplementedLocatorServer) GetLocation(context.Context, *GetLocationRequest) (*GetLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLocation not implemented")
}
func (UnimplementedLocatorServer) mustEmbedUnimplementedLocatorServer() {}
func (UnimplementedLocatorServer) testEmbeddedByValue()                 {}

// UnsafeLocatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LocatorServer will
// result in compilation errors.
type UnsafeLocatorServer interface {
	mustEmbedUnimplementedLocatorServer()
}

func RegisterLocatorServer(s grpc.ServiceRegistrar, srv LocatorServer) {
	// If the following call pancis, it indicates UnimplementedLocatorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Locator_ServiceDesc, srv)
}

func _Locator_GetLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocatorServer).GetLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Locator_GetLocation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocatorServer).GetLocation(ctx, req.(*GetLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Locator_ServiceDesc is the grpc.ServiceDesc for Locator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Locator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "locator.Locator",
	HandlerType: (*LocatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLocation",
			Handler:    _Locator_GetLocation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "goagen_locator_locator.proto",
}
