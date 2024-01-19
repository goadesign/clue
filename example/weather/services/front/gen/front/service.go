// Code generated by goa v3.14.6, DO NOT EDIT.
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
	// Endpoint for running ALL API Integration Tests for the Weather System,
	// allowing for filtering on included or excluded tests
	TestAll(context.Context, *TestAllPayload) (res *TestResults, err error)
	// Endpoint for running API Integration Tests' Smoke Tests ONLY for the Weather
	// System
	TestSmoke(context.Context) (res *TestResults, err error)
}

// APIName is the name of the API as defined in the design.
const APIName = "Weather"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "1.0.0"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "front"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [3]string{"forecast", "test_all", "test_smoke"}

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

// TestAllPayload is the payload type of the front service test_all method.
type TestAllPayload struct {
	// Tests to run
	Include []string
	// Tests to exclude
	Exclude []string
}

// Collection of test results for grouping by service
type TestCollection struct {
	// Name of the test collection
	Name string
	// Test results
	Results []*TestResult
	// Duration of the tests in ms
	Duration int64
	// Number of tests that passed
	PassCount int
	// Number of tests that failed
	FailCount int
}

// Test result for a single test
type TestResult struct {
	// Name of the test
	Name string
	// Status of the test
	Passed bool
	// Error message if the test failed
	Error *string
	// Duration of the test in ms
	Duration int64
}

// TestResults is the result type of the front service test_all method.
type TestResults struct {
	// Test collections
	Collections []*TestCollection
	// Duration of the tests in ms
	Duration int64
	// Number of tests that passed
	PassCount int
	// Number of tests that failed
	FailCount int
}

// MakeNotUsa builds a goa.ServiceError from an error.
func MakeNotUsa(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "not_usa", false, false, false)
}
