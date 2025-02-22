// Code generated by goa v3.20.0, DO NOT EDIT.
//
// Forecaster gRPC server
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/forecaster/design -o
// services/forecaster

package server

import (
	"context"

	forecaster "goa.design/clue/example/weather/services/forecaster/gen/forecaster"
	forecasterpb "goa.design/clue/example/weather/services/forecaster/gen/grpc/forecaster/pb"
	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
)

// Server implements the forecasterpb.ForecasterServer interface.
type Server struct {
	ForecastH goagrpc.UnaryHandler
	forecasterpb.UnimplementedForecasterServer
}

// New instantiates the server struct with the Forecaster service endpoints.
func New(e *forecaster.Endpoints, uh goagrpc.UnaryHandler) *Server {
	return &Server{
		ForecastH: NewForecastHandler(e.Forecast, uh),
	}
}

// NewForecastHandler creates a gRPC handler which serves the "Forecaster"
// service "forecast" endpoint.
func NewForecastHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeForecastRequest, EncodeForecastResponse)
	}
	return h
}

// Forecast implements the "Forecast" method in forecasterpb.ForecasterServer
// interface.
func (s *Server) Forecast(ctx context.Context, message *forecasterpb.ForecastRequest) (*forecasterpb.ForecastResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "forecast")
	ctx = context.WithValue(ctx, goa.ServiceKey, "Forecaster")
	resp, err := s.ForecastH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*forecasterpb.ForecastResponse), nil
}
