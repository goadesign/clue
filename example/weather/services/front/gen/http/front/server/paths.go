// Code generated by goa v3.9.1, DO NOT EDIT.
//
// HTTP request path constructors for the front service.
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/front/design -o
// services/front

package server

import (
	"fmt"
)

// ForecastFrontPath returns the URL path to the front service forecast HTTP endpoint.
func ForecastFrontPath(ip string) string {
	return fmt.Sprintf("/forecast/%v", ip)
}
