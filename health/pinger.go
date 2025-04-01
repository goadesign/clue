package health

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

	// Option configures a Pinger.
	Option func(o *options)

	client struct {
		name string
		req  *http.Request
	}

	options struct {
		scheme string
		path   string
	}
)

// NewPinger returns a new health-check client for the given service. It panics
// if the given host address is malformed.  The default scheme is "http" and the
// default path is "/livez". Both can be overridden via options.
func NewPinger(name, addr string, opts ...Option) Pinger {
	options := &options{scheme: "http", path: "/livez"}
	for _, o := range opts {
		o(options)
	}
	u := url.URL{Scheme: options.scheme, Host: addr, Path: options.path}
	req, err := http.NewRequest("GET", u.String(), nil)
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
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return fmt.Errorf("health-check for %q returned status %d", c.name, resp.StatusCode)
	}
	return nil
}

// WithScheme sets the scheme used to ping the service.
// Default scheme is "http".
func WithScheme(scheme string) Option {
	return func(o *options) {
		o.scheme = scheme
	}
}

// WithPath sets the path used to ping the service.
// Default path is "/livez".
func WithPath(path string) Option {
	return func(o *options) {
		o.path = path
	}
}
