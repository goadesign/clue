// Code generated by goa v3.5.4, DO NOT EDIT.
//
// Forecast gRPC client types
//
// Command:
// $ goa gen
// github.com/crossnokaye/micro/example/weather/services/forecast/design -o
// example/weather/services/forecast

package client

import (
	forecast "github.com/crossnokaye/micro/example/weather/services/forecast/gen/forecast"
	forecastpb "github.com/crossnokaye/micro/example/weather/services/forecast/gen/grpc/forecast/pb"
	goa "goa.design/goa/v3/pkg"
)

// NewForecastRequest builds the gRPC request type from the payload of the
// "forecast" endpoint of the "Forecast" service.
func NewForecastRequest(payload *forecast.ForecastPayload) *forecastpb.ForecastRequest {
	message := &forecastpb.ForecastRequest{
		Lat:  payload.Lat,
		Long: payload.Long,
	}
	return message
}

// NewForecastResult builds the result type of the "forecast" endpoint of the
// "Forecast" service from the gRPC response type.
func NewForecastResult(message *forecastpb.ForecastResponse) *forecast.Forecast2 {
	result := &forecast.Forecast2{}
	if message.Location != nil {
		result.Location = protobufForecastpbLocationToForecastLocation(message.Location)
	}
	if message.Periods != nil {
		result.Periods = make([]*forecast.Period, len(message.Periods))
		for i, val := range message.Periods {
			result.Periods[i] = &forecast.Period{
				Name:            val.Name,
				StartTime:       val.StartTime,
				EndTime:         val.EndTime,
				Temperature:     int(val.Temperature),
				TemperatureUnit: val.TemperatureUnit,
				Summary:         val.Summary,
			}
		}
	}
	return result
}

// ValidateForecastResponse runs the validations defined on ForecastResponse.
func ValidateForecastResponse(message *forecastpb.ForecastResponse) (err error) {
	if message.Location == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("location", "message"))
	}
	if message.Periods == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("periods", "message"))
	}
	for _, e := range message.Periods {
		if e != nil {
			if err2 := ValidatePeriod(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateLocation runs the validations defined on Location.
func ValidateLocation(message *forecastpb.Location) (err error) {

	return
}

// ValidatePeriod runs the validations defined on Period.
func ValidatePeriod(message *forecastpb.Period) (err error) {
	err = goa.MergeErrors(err, goa.ValidateFormat("message.startTime", message.StartTime, goa.FormatDateTime))

	err = goa.MergeErrors(err, goa.ValidateFormat("message.endTime", message.EndTime, goa.FormatDateTime))

	return
}

// svcForecastLocationToForecastpbLocation builds a value of type
// *forecastpb.Location from a value of type *forecast.Location.
func svcForecastLocationToForecastpbLocation(v *forecast.Location) *forecastpb.Location {
	res := &forecastpb.Location{
		Lat:   v.Lat,
		Long:  v.Long,
		City:  v.City,
		State: v.State,
	}

	return res
}

// protobufForecastpbLocationToForecastLocation builds a value of type
// *forecast.Location from a value of type *forecastpb.Location.
func protobufForecastpbLocationToForecastLocation(v *forecastpb.Location) *forecast.Location {
	res := &forecast.Location{
		Lat:   v.Lat,
		Long:  v.Long,
		City:  v.City,
		State: v.State,
	}

	return res
}
