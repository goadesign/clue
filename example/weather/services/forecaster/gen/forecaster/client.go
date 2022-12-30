// Code generated by goa v3.10.2, DO NOT EDIT.
//
// Forecaster client
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/forecaster/design -o
// services/forecaster

package forecaster

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "Forecaster" service client.
type Client struct {
	ForecastEndpoint goa.Endpoint
}

// NewClient initializes a "Forecaster" service client given the endpoints.
func NewClient(forecast goa.Endpoint) *Client {
	return &Client{
		ForecastEndpoint: forecast,
	}
}

// Forecast calls the "forecast" endpoint of the "Forecaster" service.
func (c *Client) Forecast(ctx context.Context, p *ForecastPayload) (res *Forecast2, err error) {
	var ires interface{}
	ires, err = c.ForecastEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Forecast2), nil
}
