// Code generated by goa v3.5.4, DO NOT EDIT.
//
// HTTP request path constructors for the front service.
//
// Command:
// $ goa gen github.com/crossnokaye/micro/example/weather/services/front/design
// -o example/weather/services/front

package client

import (
	"fmt"
)

// ForecastFrontPath returns the URL path to the front service forecast HTTP endpoint.
func ForecastFrontPath(ip string) string {
	return fmt.Sprintf("/forecast/%v", ip)
}
