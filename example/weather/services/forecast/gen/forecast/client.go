// Code generated by goa v3.5.4, DO NOT EDIT.
//
// Forecast client
//
// Command:
// $ goa gen
// github.com/crossnokaye/micro/example/weather/services/forecast/design -o
// services/forecast

package forecast

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "Forecast" service client.
type Client struct {
	ForecastEndpoint goa.Endpoint
}

// NewClient initializes a "Forecast" service client given the endpoints.
func NewClient(forecast goa.Endpoint) *Client {
	return &Client{
		ForecastEndpoint: forecast,
	}
}

// Forecast calls the "forecast" endpoint of the "Forecast" service.
func (c *Client) Forecast(ctx context.Context, p *ForecastPayload) (res *Forecast2, err error) {
	var ires interface{}
	ires, err = c.ForecastEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Forecast2), nil
}