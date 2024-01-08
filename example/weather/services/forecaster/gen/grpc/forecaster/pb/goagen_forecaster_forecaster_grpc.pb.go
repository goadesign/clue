// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.0
// source: goagen_forecaster_forecaster.proto

package forecasterpb

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

// ForecasterClient is the client API for Forecaster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ForecasterClient interface {
	// Retrieve weather forecast for a given location
	Forecast(ctx context.Context, in *ForecastRequest, opts ...grpc.CallOption) (*ForecastResponse, error)
}

type forecasterClient struct {
	cc grpc.ClientConnInterface
}

func NewForecasterClient(cc grpc.ClientConnInterface) ForecasterClient {
	return &forecasterClient{cc}
}

func (c *forecasterClient) Forecast(ctx context.Context, in *ForecastRequest, opts ...grpc.CallOption) (*ForecastResponse, error) {
	out := new(ForecastResponse)
	err := c.cc.Invoke(ctx, "/forecaster.Forecaster/Forecast", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ForecasterServer is the server API for Forecaster service.
// All implementations must embed UnimplementedForecasterServer
// for forward compatibility
type ForecasterServer interface {
	// Retrieve weather forecast for a given location
	Forecast(context.Context, *ForecastRequest) (*ForecastResponse, error)
	mustEmbedUnimplementedForecasterServer()
}

// UnimplementedForecasterServer must be embedded to have forward compatible implementations.
type UnimplementedForecasterServer struct {
}

func (UnimplementedForecasterServer) Forecast(context.Context, *ForecastRequest) (*ForecastResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Forecast not implemented")
}
func (UnimplementedForecasterServer) mustEmbedUnimplementedForecasterServer() {}

// UnsafeForecasterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ForecasterServer will
// result in compilation errors.
type UnsafeForecasterServer interface {
	mustEmbedUnimplementedForecasterServer()
}

func RegisterForecasterServer(s grpc.ServiceRegistrar, srv ForecasterServer) {
	s.RegisterService(&Forecaster_ServiceDesc, srv)
}

func _Forecaster_Forecast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForecastRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForecasterServer).Forecast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forecaster.Forecaster/Forecast",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForecasterServer).Forecast(ctx, req.(*ForecastRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Forecaster_ServiceDesc is the grpc.ServiceDesc for Forecaster service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Forecaster_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "forecaster.Forecaster",
	HandlerType: (*ForecasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Forecast",
			Handler:    _Forecaster_Forecast_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "goagen_forecaster_forecaster.proto",
}