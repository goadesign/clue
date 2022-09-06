package log

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"goa.design/goa/v3/middleware"
)

type (
	// Doer is a client that can execute HTTP requests.
	Doer interface {
		Do(*http.Request) (*http.Response, error)
	}

	// HTTPLogOption is a function that applies a configuration option to log
	// HTTP middleware.
	HTTPLogOption func(*httpLogOptions)

	// HTTPClientLogOption is a function that applies a configuration option
	// to a HTTP client logger.
	HTTPClientLogOption func(*httpClientOptions)

	httpLogOptions struct {
		pathFilters []*regexp.Regexp
	}

	httpClientOptions struct {
		iserr func(int) bool
	}
	// client wraps an HTTP client and logs requests and responses.
	client struct {
		Doer
		options *httpClientOptions
	}
)

// HTTP returns a HTTP middleware that initializes the logger context.  It
// panics if logCtx was not initialized with Context.
func HTTP(logCtx context.Context, opts ...HTTPLogOption) func(http.Handler) http.Handler {
	MustContainLogger(logCtx)
	var options httpLogOptions
	for _, o := range opts {
		o(&options)
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			for _, opt := range options.pathFilters {
				if opt.MatchString(req.URL.Path) {
					h.ServeHTTP(w, req)
					return
				}
			}
			ctx := WithContext(req.Context(), logCtx)
			if requestID := req.Context().Value(middleware.RequestIDKey); requestID != nil {
				ctx = With(ctx, KV{RequestIDKey, requestID})
			}
			h.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// Client returns a HTTP client that wraps the given doer and log requests and
// responses using the clue logger stored in the request context.
func Client(doer Doer, opts ...HTTPClientLogOption) Doer {
	options := &httpClientOptions{
		iserr: func(status int) bool { return status >= 400 },
	}
	for _, o := range opts {
		o(options)
	}
	return &client{Doer: doer, options: options}
}

// WithPathFilter adds a path filter to the HTTP middleware. Requests whose path
// match the filter are not logged. WithPathFilter can be called multiple times
// to add multiple filters.
func WithPathFilter(filter *regexp.Regexp) HTTPLogOption {
	return func(o *httpLogOptions) {
		o.pathFilters = append(o.pathFilters, filter)
	}
}

// WithErrorStatus returns a HTTP client logger option that configures the
// logger to log errors for responses with the given status code.
func WithErrorStatus(status int) HTTPClientLogOption {
	return func(o *httpClientOptions) {
		o.iserr = func(s int) bool { return s == status }
	}
}

// Do executes the given HTTP request and logs the request and response. The
// request context must be initialized with a clue logger.
func (c *client) Do(req *http.Request) (resp *http.Response, err error) {
	msgKV := KV{K: MessageKey, V: "finished client HTTP request"}
	methKV := KV{K: HTTPMethodKey, V: req.Method}
	urlKV := KV{K: HTTPURLKey, V: req.URL.String()}
	then := timeNow()
	resp, err = c.Doer.Do(req)
	if err != nil {
		Error(req.Context(), err, msgKV, methKV, urlKV)
		return
	}
	ms := timeSince(then).Milliseconds()
	statusKV := KV{K: HTTPStatusKey, V: resp.Status}
	durKV := KV{K: HTTPDurationKey, V: ms}
	if c.options.iserr(resp.StatusCode) {
		Error(req.Context(), fmt.Errorf(resp.Status), msgKV, methKV, urlKV, statusKV, durKV)
		return
	}
	Print(req.Context(), msgKV, methKV, urlKV, statusKV, durKV)
	return
}
