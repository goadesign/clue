// Code generated by goa v3.9.1, DO NOT EDIT.
//
// test endpoints
//
// Command:
// $ goa gen goa.design/clue/internal/testsvc/design

package test

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "test" service endpoints.
type Endpoints struct {
	HTTPMethod goa.Endpoint
	GrpcMethod goa.Endpoint
	GrpcStream goa.Endpoint
}

// GrpcStreamEndpointInput holds both the payload and the server stream of the
// "grpc_stream" method.
type GrpcStreamEndpointInput struct {
	// Stream is the server stream used by the "grpc_stream" method to send data.
	Stream GrpcStreamServerStream
}

// NewEndpoints wraps the methods of the "test" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		HTTPMethod: NewHTTPMethodEndpoint(s),
		GrpcMethod: NewGrpcMethodEndpoint(s),
		GrpcStream: NewGrpcStreamEndpoint(s),
	}
}

// Use applies the given middleware to all the "test" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.HTTPMethod = m(e.HTTPMethod)
	e.GrpcMethod = m(e.GrpcMethod)
	e.GrpcStream = m(e.GrpcStream)
}

// NewHTTPMethodEndpoint returns an endpoint function that calls the method
// "http_method" of service "test".
func NewHTTPMethodEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*Fields)
		return s.HTTPMethod(ctx, p)
	}
}

// NewGrpcMethodEndpoint returns an endpoint function that calls the method
// "grpc_method" of service "test".
func NewGrpcMethodEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*Fields)
		return s.GrpcMethod(ctx, p)
	}
}

// NewGrpcStreamEndpoint returns an endpoint function that calls the method
// "grpc_stream" of service "test".
func NewGrpcStreamEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		ep := req.(*GrpcStreamEndpointInput)
		return nil, s.GrpcStream(ctx, ep.Stream)
	}
}
