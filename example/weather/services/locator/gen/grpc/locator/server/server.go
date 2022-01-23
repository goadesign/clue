// Code generated by goa v3.5.4, DO NOT EDIT.
//
// locator gRPC server
//
// Command:
// $ goa gen
// github.com/crossnokaye/micro/example/weather/services/locator/design -o
// services/locator

package server

import (
	"context"

	locatorpb "github.com/crossnokaye/micro/example/weather/services/locator/gen/grpc/locator/pb"
	locator "github.com/crossnokaye/micro/example/weather/services/locator/gen/locator"
	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
)

// Server implements the locatorpb.LocatorServer interface.
type Server struct {
	GetLocationH goagrpc.UnaryHandler
	locatorpb.UnimplementedLocatorServer
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the expr.
type ErrorNamer interface {
	ErrorName() string
}

// New instantiates the server struct with the locator service endpoints.
func New(e *locator.Endpoints, uh goagrpc.UnaryHandler) *Server {
	return &Server{
		GetLocationH: NewGetLocationHandler(e.GetLocation, uh),
	}
}

// NewGetLocationHandler creates a gRPC handler which serves the "locator"
// service "get_location" endpoint.
func NewGetLocationHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeGetLocationRequest, EncodeGetLocationResponse)
	}
	return h
}

// GetLocation implements the "GetLocation" method in locatorpb.LocatorServer
// interface.
func (s *Server) GetLocation(ctx context.Context, message *locatorpb.GetLocationRequest) (*locatorpb.GetLocationResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "get_location")
	ctx = context.WithValue(ctx, goa.ServiceKey, "locator")
	resp, err := s.GetLocationH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*locatorpb.GetLocationResponse), nil
}
