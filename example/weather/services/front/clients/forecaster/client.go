package forecaster

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"

	genforecast "goa.design/clue/example/weather/services/forecaster/gen/forecaster"
	genclient "goa.design/clue/example/weather/services/forecaster/gen/grpc/forecaster/client"
)

type (
	// Client is a client for the forecast service.
	Client interface {
		// GetForecast gets the forecast for the given location.
		GetForecast(ctx context.Context, lat, long float64) (*Forecast, error)
	}

	// Forecast represents the forecast for a given location.
	Forecast struct {
		// Location is the location of the forecast.
		Location *Location
		// Periods is the forecast for the location.
		Periods []*Period
	}

	// Location represents the geographical location of a forecast.
	Location struct {
		// Lat is the latitude of the location.
		Lat float64
		// Long is the longitude of the location.
		Long float64
		// City is the city of the location.
		City string
		// State is the state of the location.
		State string
	}

	// Period represents a forecast period.
	Period struct {
		// Name is the name of the forecast period.
		Name string
		// StartTime is the start time of the forecast period in RFC3339 format.
		StartTime string
		// EndTime is the end time of the forecast period in RFC3339 format.
		EndTime string
		// Temperature is the temperature of the forecast period.
		Temperature int
		// TemperatureUnit is the temperature unit of the forecast period.
		TemperatureUnit string
		// Summary is the summary of the forecast period.
		Summary string
	}

	// client is the client implementation.
	client struct {
		forecast goa.Endpoint
	}
)

// New instantiates a new forecast service client.
func New(cc *grpc.ClientConn) Client {
	c := genclient.NewClient(cc, grpc.WaitForReady(true))
	return &client{c.Forecast()}
}

// Forecast returns the forecast for the given location or current location if
// lat or long are nil.
func (c *client) GetForecast(ctx context.Context, lat, long float64) (*Forecast, error) {
	res, err := c.forecast(ctx, &genforecast.ForecastPayload{lat, long})
	if err != nil {
		return nil, err
	}
	f := res.(*genforecast.Forecast2)
	l := Location(*f.Location)
	ps := make([]*Period, len(f.Periods))
	for i, p := range f.Periods {
		pval := Period(*p)
		ps[i] = &pval
	}
	return &Forecast{&l, ps}, nil
}
