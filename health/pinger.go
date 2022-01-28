package health

import (
	"context"
	"fmt"
	"net/http"
)

type (
	// Pinger makes it possible to ping a service.
	Pinger interface {
		// Name of remote service.
		Name() string
		// Ping the remote service, return a non nil error if the
		// service is not available.
		Ping(context.Context) error
	}

	client struct {
		name string
		req  *http.Request
	}
)

// NewPinger returns a new health-check client for the given service.
// NewPinger panics if the given host address is malformed.
func NewPinger(name, scheme, addr string) Pinger {
	if scheme == "" {
		scheme = "http"
	}
	url := scheme + "://" + addr + "/livez"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	return &client{
		name: name,
		req:  req,
	}
}

func (c *client) Name() string {
	return c.name
}

func (c *client) Ping(ctx context.Context) error {
	resp, err := http.DefaultClient.Do(c.req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to make health check request to %q: %v", c.name, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return fmt.Errorf("health-check for %q returned status %d", c.name, resp.StatusCode)
	}
	return nil
}
