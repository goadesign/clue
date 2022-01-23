package ipapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type (
	// Client is the ipapi client interface. It implements a simple
	// in-memory transparent cache to avoid hitting the API too often and
	// getting throttled (429).
	Client interface {
		// GetLocation gets the location for the given IP address or the
		// IP of the server if blank.
		GetLocation(ctx context.Context, ip string) (*WorldLocation, error)
		// Name provides a client name used to report health check issues.
		Name() string
		// Ping checks the client is healthy.
		Ping(ctx context.Context) error
	}

	// WorldLocation represents the geographical location of an IP address.
	WorldLocation struct {
		// Lat is the latitude of the location.
		Lat float64 `json:"lat"`
		// Long is the longitude of the location.
		Long float64 `json:"lon"`
		// City is the city of the location.
		City string `json:"city"`
		// Region is the region/state of the location.
		Region string `json:"region"`
		// Country is the country of the location.
		Country string `json:"country"`
	}

	// client implements Client.
	client struct {
		c     *http.Client
		cache map[string]*WorldLocation
	}
)

// baseURL is the base URL for the ipapi service.
const baseURL = "http://ip-api.com/json"

// New returns a new client for the ipapi.co API.
func New(c *http.Client) Client {
	return &client{c: c, cache: make(map[string]*WorldLocation)}
}

// GetLocation gets the location for the given IP address.
func (c *client) GetLocation(ctx context.Context, ip string) (*WorldLocation, error) {
	if l, ok := c.cache[ip]; ok {
		return l, nil
	}
	body, err := c.getLocation(ctx, ip)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	var l WorldLocation
	if err := json.NewDecoder(body).Decode(&l); err != nil {
		return nil, err
	}
	c.cache[ip] = &l
	return &l, nil
}

// Name provides a client name used to report health check issues.
func (c *client) Name() string {
	return "ip-api"
}

// Ping checks the client is healthy.
func (c *client) Ping(ctx context.Context) error {
	body, err := c.getLocation(ctx, "")
	if err != nil {
		return err
	}
	body.Close()
	return nil
}

// getLocation gets the location for the given IP address.
func (c *client) getLocation(ctx context.Context, ip string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", baseURL, ip), nil)
	if err != nil {
		return nil, fmt.Errorf("invalid IP address %q", ip)
	}
	resp, err := c.c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			msg = []byte("unknown error")
		}
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d (%s)", resp.StatusCode, string(msg))
	}
	return resp.Body, nil
}
