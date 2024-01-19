// Code generated by goa v3.14.6, DO NOT EDIT.
//
// tester gRPC server
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/tester/design -o
// services/tester

package server

import (
	"context"
	"errors"

	testerpb "goa.design/clue/example/weather/services/tester/gen/grpc/tester/pb"
	tester "goa.design/clue/example/weather/services/tester/gen/tester"
	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc/codes"
)

// Server implements the testerpb.TesterServer interface.
type Server struct {
	TestAllH        goagrpc.UnaryHandler
	TestSmokeH      goagrpc.UnaryHandler
	TestForecasterH goagrpc.UnaryHandler
	TestLocatorH    goagrpc.UnaryHandler
	testerpb.UnimplementedTesterServer
}

// New instantiates the server struct with the tester service endpoints.
func New(e *tester.Endpoints, uh goagrpc.UnaryHandler) *Server {
	return &Server{
		TestAllH:        NewTestAllHandler(e.TestAll, uh),
		TestSmokeH:      NewTestSmokeHandler(e.TestSmoke, uh),
		TestForecasterH: NewTestForecasterHandler(e.TestForecaster, uh),
		TestLocatorH:    NewTestLocatorHandler(e.TestLocator, uh),
	}
}

// NewTestAllHandler creates a gRPC handler which serves the "tester" service
// "test_all" endpoint.
func NewTestAllHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeTestAllRequest, EncodeTestAllResponse)
	}
	return h
}

// TestAll implements the "TestAll" method in testerpb.TesterServer interface.
func (s *Server) TestAll(ctx context.Context, message *testerpb.TestAllRequest) (*testerpb.TestAllResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "test_all")
	ctx = context.WithValue(ctx, goa.ServiceKey, "tester")
	resp, err := s.TestAllH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "include_exclude_both":
				return nil, goagrpc.NewStatusError(codes.InvalidArgument, err, goagrpc.NewErrorResponse(err))
			case "wildcard_compile_error":
				return nil, goagrpc.NewStatusError(codes.InvalidArgument, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*testerpb.TestAllResponse), nil
}

// NewTestSmokeHandler creates a gRPC handler which serves the "tester" service
// "test_smoke" endpoint.
func NewTestSmokeHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, nil, EncodeTestSmokeResponse)
	}
	return h
}

// TestSmoke implements the "TestSmoke" method in testerpb.TesterServer
// interface.
func (s *Server) TestSmoke(ctx context.Context, message *testerpb.TestSmokeRequest) (*testerpb.TestSmokeResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "test_smoke")
	ctx = context.WithValue(ctx, goa.ServiceKey, "tester")
	resp, err := s.TestSmokeH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*testerpb.TestSmokeResponse), nil
}

// NewTestForecasterHandler creates a gRPC handler which serves the "tester"
// service "test_forecaster" endpoint.
func NewTestForecasterHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, nil, EncodeTestForecasterResponse)
	}
	return h
}

// TestForecaster implements the "TestForecaster" method in
// testerpb.TesterServer interface.
func (s *Server) TestForecaster(ctx context.Context, message *testerpb.TestForecasterRequest) (*testerpb.TestForecasterResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "test_forecaster")
	ctx = context.WithValue(ctx, goa.ServiceKey, "tester")
	resp, err := s.TestForecasterH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "include_exclude_both":
				return nil, goagrpc.NewStatusError(codes.InvalidArgument, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*testerpb.TestForecasterResponse), nil
}

// NewTestLocatorHandler creates a gRPC handler which serves the "tester"
// service "test_locator" endpoint.
func NewTestLocatorHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, nil, EncodeTestLocatorResponse)
	}
	return h
}

// TestLocator implements the "TestLocator" method in testerpb.TesterServer
// interface.
func (s *Server) TestLocator(ctx context.Context, message *testerpb.TestLocatorRequest) (*testerpb.TestLocatorResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "test_locator")
	ctx = context.WithValue(ctx, goa.ServiceKey, "tester")
	resp, err := s.TestLocatorH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "include_exclude_both":
				return nil, goagrpc.NewStatusError(codes.InvalidArgument, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*testerpb.TestLocatorResponse), nil
}
