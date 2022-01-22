// Code generated by goa v3.5.4, DO NOT EDIT.
//
// Forecast gRPC server encoders and decoders
//
// Command:
// $ goa gen
// github.com/crossnokaye/micro/example/weather/services/forecast/design -o
// example/weather/services/forecast

package server

import (
	"context"

	forecast "github.com/crossnokaye/micro/example/weather/services/forecast/gen/forecast"
	forecastpb "github.com/crossnokaye/micro/example/weather/services/forecast/gen/grpc/forecast/pb"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc/metadata"
)

// EncodeForecastResponse encodes responses from the "Forecast" service
// "forecast" endpoint.
func EncodeForecastResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.(*forecast.Forecast2)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Forecast", "forecast", "*forecast.Forecast2", v)
	}
	resp := NewForecastResponse(result)
	return resp, nil
}

// DecodeForecastRequest decodes requests sent to "Forecast" service "forecast"
// endpoint.
func DecodeForecastRequest(ctx context.Context, v interface{}, md metadata.MD) (interface{}, error) {
	var (
		message *forecastpb.ForecastRequest
		ok      bool
	)
	{
		if message, ok = v.(*forecastpb.ForecastRequest); !ok {
			return nil, goagrpc.ErrInvalidType("Forecast", "forecast", "*forecastpb.ForecastRequest", v)
		}
	}
	var payload *forecast.ForecastPayload
	{
		payload = NewForecastPayload(message)
	}
	return payload, nil
}
