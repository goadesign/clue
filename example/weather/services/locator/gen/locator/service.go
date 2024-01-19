// Code generated by goa v3.14.6, DO NOT EDIT.
//
// locator service
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/locator/design -o
// services/locator

package locator

import (
	"context"
)

// Public HTTP frontend
type Service interface {
	// Retrieve location information for a given IP address
	GetLocation(context.Context, string) (res *WorldLocation, err error)
}

// APIName is the name of the API as defined in the design.
const APIName = "IP Location API"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "1.0.0"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "locator"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"get_location"}

// WorldLocation is the result type of the locator service get_location method.
type WorldLocation struct {
	// Latitude
	Lat float64
	// Longitude
	Long float64
	// City
	City string
	// State, region etc.
	Region string
	// Country
	Country string
}
