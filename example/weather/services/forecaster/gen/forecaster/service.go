// Code generated by goa v3.5.4, DO NOT EDIT.
//
// Forecaster service
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/forecaster/design -o
// services/forecaster

package forecaster

import (
	"context"
)

// Service that provides weather forecasts
type Service interface {
	// Retrieve weather forecast for a given location
	Forecast(context.Context, *ForecastPayload) (res *Forecast2, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "Forecaster"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"forecast"}

// ForecastPayload is the payload type of the Forecaster service forecast
// method.
type ForecastPayload struct {
	// Latitude
	Lat float64
	// Longitude
	Long float64
}

// Forecast2 is the result type of the Forecaster service forecast method.
type Forecast2 struct {
	// Forecast location
	Location *Location
	// Weather forecast periods
	Periods []*Period
}

// Geographical location
type Location struct {
	// Latitude
	Lat float64
	// Longitude
	Long float64
	// City
	City string
	// State
	State string
}

// Weather forecast period
type Period struct {
	// Period name
	Name string
	// Start time
	StartTime string
	// End time
	EndTime string
	// Temperature
	Temperature int
	// Temperature unit
	TemperatureUnit string
	// Summary
	Summary string
}
