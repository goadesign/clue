// Code generated by goa v3.20.0, DO NOT EDIT.
//
// Forecaster gRPC client types
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/forecaster/design -o
// services/forecaster

package client

import (
	forecaster "goa.design/clue/example/weather/services/forecaster/gen/forecaster"
	forecasterpb "goa.design/clue/example/weather/services/forecaster/gen/grpc/forecaster/pb"
	goa "goa.design/goa/v3/pkg"
)

// NewProtoForecastRequest builds the gRPC request type from the payload of the
// "forecast" endpoint of the "Forecaster" service.
func NewProtoForecastRequest(payload *forecaster.ForecastPayload) *forecasterpb.ForecastRequest {
	message := &forecasterpb.ForecastRequest{
		Lat:  payload.Lat,
		Long: payload.Long,
	}
	return message
}

// NewForecastResult builds the result type of the "forecast" endpoint of the
// "Forecaster" service from the gRPC response type.
func NewForecastResult(message *forecasterpb.ForecastResponse) *forecaster.Forecast2 {
	result := &forecaster.Forecast2{}
	if message.Location != nil {
		result.Location = protobufForecasterpbLocationToForecasterLocation(message.Location)
	}
	if message.Periods != nil {
		result.Periods = make([]*forecaster.Period, len(message.Periods))
		for i, val := range message.Periods {
			result.Periods[i] = &forecaster.Period{
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
func ValidateForecastResponse(message *forecasterpb.ForecastResponse) (err error) {
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

// ValidatePeriod runs the validations defined on Period.
func ValidatePeriod(elem *forecasterpb.Period) (err error) {
	err = goa.MergeErrors(err, goa.ValidateFormat("elem.startTime", elem.StartTime, goa.FormatDateTime))
	err = goa.MergeErrors(err, goa.ValidateFormat("elem.endTime", elem.EndTime, goa.FormatDateTime))
	return
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
