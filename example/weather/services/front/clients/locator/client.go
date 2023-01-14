package locator

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"

	"goa.design/clue/debug"
	genclient "goa.design/clue/example/weather/services/locator/gen/grpc/locator/client"
	genlocator "goa.design/clue/example/weather/services/locator/gen/locator"
)

type (
	// Client is a client for the locator service.
	Client interface {
		// GetLocation gets the location for the given ip
		GetLocation(ctx context.Context, ip string) (*WorldLocation, error)
	}

	// WorldLocation represents the location for the given IP address.
	WorldLocation struct {
		// Lat is the latitude of the location.
		Lat float64
		// Long is the longitude of the location.
		Long float64
		// City is the city of the location.
		City string
		// Region is the region/state of the location.
		Region string
		// Country is the country of the location.
		Country string
	}

	// client is the client implementation.
	client struct {
		locator goa.Endpoint
	}
)

// New instantiates a new locator service client.
func New(cc *grpc.ClientConn) Client {
	c := genclient.NewClient(cc, grpc.WaitForReady(true))
	return &client{debug.LogPayloads(debug.WithClient())(c.GetLocation())}
}

// GetLocation returns the location for the given IP address.
func (c *client) GetLocation(ctx context.Context, ip string) (*WorldLocation, error) {
	res, err := c.locator(ctx, ip)
	if err != nil {
		return nil, err
	}
	ipl := res.(*genlocator.WorldLocation)
	l := WorldLocation(*ipl)
	return &l, nil
}
