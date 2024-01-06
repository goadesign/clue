// Code generated by goa v3.14.2, DO NOT EDIT.
//
// HTTP request path constructors for the front service.
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/front/design -o
// services/front

package client

import (
	"fmt"
)

// ForecastFrontPath returns the URL path to the front service forecast HTTP endpoint.
func ForecastFrontPath(ip string) string {
	return fmt.Sprintf("/forecast/%v", ip)
}

// TestAllFrontPath returns the URL path to the front service test_all HTTP endpoint.
func TestAllFrontPath() string {
	return "/tester/all"
}

// TestSmokeFrontPath returns the URL path to the front service test_smoke HTTP endpoint.
func TestSmokeFrontPath() string {
	return "/tester/smoke"
}
