// Code generated by goa v3.20.0, DO NOT EDIT.
//
// test service
//
// Command:
// $ goa gen goa.design/clue/internal/testsvc/design

package test

import (
	"context"
)

// Service is the test service interface.
type Service interface {
	// HTTPMethod implements http_method.
	HTTPMethod(context.Context, *Fields) (res *Fields, err error)
	// GrpcMethod implements grpc_method.
	GrpcMethod(context.Context, *Fields) (res *Fields, err error)
	// GrpcStream implements grpc_stream.
	GrpcStream(context.Context, GrpcStreamServerStream) (err error)
}

// APIName is the name of the API as defined in the design.
const APIName = "itest"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "0.0.1"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "test"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [3]string{"http_method", "grpc_method", "grpc_stream"}

// GrpcStreamServerStream is the interface a "grpc_stream" endpoint server
// stream must satisfy.
type GrpcStreamServerStream interface {
	// Send streams instances of "Fields".
	Send(*Fields) error
	// SendWithContext streams instances of "Fields" with context.
	SendWithContext(context.Context, *Fields) error
	// Recv reads instances of "Fields" from the stream.
	Recv() (*Fields, error)
	// RecvWithContext reads instances of "Fields" from the stream with context.
	RecvWithContext(context.Context) (*Fields, error)
	// Close closes the stream.
	Close() error
}

// GrpcStreamClientStream is the interface a "grpc_stream" endpoint client
// stream must satisfy.
type GrpcStreamClientStream interface {
	// Send streams instances of "Fields".
	Send(*Fields) error
	// SendWithContext streams instances of "Fields" with context.
	SendWithContext(context.Context, *Fields) error
	// Recv reads instances of "Fields" from the stream.
	Recv() (*Fields, error)
	// RecvWithContext reads instances of "Fields" from the stream with context.
	RecvWithContext(context.Context) (*Fields, error)
	// Close closes the stream.
	Close() error
}

// Fields is the payload type of the test service http_method method.
type Fields struct {
	// String operand
	S *string
	// Int operand
	I *int
}
