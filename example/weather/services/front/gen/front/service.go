// Code generated by goa v3.9.1, DO NOT EDIT.
//
// front service
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/front/design -o
// services/front

package front

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Public HTTP frontend
type Service interface {
	// Retrieve weather forecast for given IP
	Forecast(context.Context, string) (res *Forecast2, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "front"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"forecast"}

// Forecast2 is the result type of the front service forecast method.
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

// MakeNotUsa builds a goa.ServiceError from an error.
func MakeNotUsa(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "not_usa", false, false, false)
}
