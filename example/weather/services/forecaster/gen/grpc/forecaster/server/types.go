// Code generated by goa v3.8.3, DO NOT EDIT.
//
// Forecaster gRPC server types
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/forecaster/design -o
// services/forecaster

package server

import (
	forecaster "goa.design/clue/example/weather/services/forecaster/gen/forecaster"
	forecasterpb "goa.design/clue/example/weather/services/forecaster/gen/grpc/forecaster/pb"
)

// NewForecastPayload builds the payload of the "forecast" endpoint of the
// "Forecaster" service from the gRPC request type.
func NewForecastPayload(message *forecasterpb.ForecastRequest) *forecaster.ForecastPayload {
	v := &forecaster.ForecastPayload{
		Lat:  message.Lat,
		Long: message.Long,
	}
	return v
}

// NewProtoForecastResponse builds the gRPC response type from the result of
// the "forecast" endpoint of the "Forecaster" service.
func NewProtoForecastResponse(result *forecaster.Forecast2) *forecasterpb.ForecastResponse {
	message := &forecasterpb.ForecastResponse{}
	if result.Location != nil {
		message.Location = svcForecasterLocationToForecasterpbLocation(result.Location)
	}
	if result.Periods != nil {
		message.Periods = make([]*forecasterpb.Period, len(result.Periods))
		for i, val := range result.Periods {
			message.Periods[i] = &forecasterpb.Period{
				Name:            val.Name,
				StartTime:       val.StartTime,
				EndTime:         val.EndTime,
				Temperature:     int32(val.Temperature),
				TemperatureUnit: val.TemperatureUnit,
				Summary:         val.Summary,
			}
		}
	}
	return message
}

// svcForecasterLocationToForecasterpbLocation builds a value of type
// *forecasterpb.Location from a value of type *forecaster.Location.
func svcForecasterLocationToForecasterpbLocation(v *forecaster.Location) *forecasterpb.Location {
	res := &forecasterpb.Location{
		Lat:   v.Lat,
		Long:  v.Long,
		City:  v.City,
		State: v.State,
	}

	return res
}

// protobufForecasterpbLocationToForecasterLocation builds a value of type
// *forecaster.Location from a value of type *forecasterpb.Location.
func protobufForecasterpbLocationToForecasterLocation(v *forecasterpb.Location) *forecaster.Location {
	res := &forecaster.Location{
		Lat:   v.Lat,
		Long:  v.Long,
		City:  v.City,
		State: v.State,
	}

	return res
}
