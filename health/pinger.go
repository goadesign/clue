package health

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
		name       string
		httpClient *http.Client
		httpURL    string
	}

	options struct {
		scheme    string
		path      string
		timeout   time.Duration
		transport http.RoundTripper
	}
)

// NewPinger returns a new health-check client for the given service. It panics
// if the given host address is malformed.  The default scheme is "http" and the
// default path is "/livez". Both can be overridden via options.
func NewPinger(name, addr string, opts ...Option) Pinger {
	options := &options{scheme: "http", path: "/livez", transport: otelhttp.NewTransport(http.DefaultTransport)}
	for _, o := range opts {
		o(options)
	}
	u := url.URL{Scheme: options.scheme, Host: addr, Path: options.path}

	return &client{
		name:       name,
		httpClient: &http.Client{Timeout: options.timeout, Transport: options.transport},
		httpURL:    u.String(),
	}
}

func (c *client) Name() string {
	return c.name
}

func (c *client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.httpURL, nil)
	if err != nil {
		return fmt.Errorf("failed to prepare health check request to %q: %v", c.name, err)
	}
	resp, err := c.httpClient.Do(req)
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

func WithTransport(transport http.RoundTripper) Option {
	return func(o *options) {
		o.transport = transport
	}
}

// WithTimeout sets the timeout used to ping the service.
// Default is no timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}
