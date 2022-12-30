package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"goa.design/goa/v3/middleware"
	goa "goa.design/goa/v3/pkg"
)

type (
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
		iserr      func(int) bool
		logErrBody bool
	}

	// client wraps an HTTP roundtripper and logs requests and responses.
	client struct {
		http.RoundTripper
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

// Endpoint is a Goa endpoint middleware that adds the service and method names
// to the logged key/value pairs.
func Endpoint(e goa.Endpoint) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if s := ctx.Value(goa.ServiceKey); s != nil {
			ctx = With(ctx, KV{K: "goa.service", V: s})
		}
		if m := ctx.Value(goa.MethodKey); m != nil {
			ctx = With(ctx, KV{K: "goa.method", V: m})
		}
		return e(ctx, req)
	}
}

// Client wraps the given roundtripper and log requests and responses using the
// clue logger stored in the request context.
func Client(t http.RoundTripper, opts ...HTTPClientLogOption) http.RoundTripper {
	options := &httpClientOptions{
		iserr: func(status int) bool { return status >= 400 },
	}
	for _, o := range opts {
		o(options)
	}
	return &client{RoundTripper: t, options: options}
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

// WithLogBodyOnError returns a HTTP client logger option that configures the
// logger to log the response body when the response status code is an error.
func WithLogBodyOnError() HTTPClientLogOption {
	return func(o *httpClientOptions) {
		o.logErrBody = true
	}
}

// RoundTrip executes the given HTTP request and logs the request and response. The
// request context must be initialized with a clue logger.
func (c *client) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	msgKV := KV{K: MessageKey, V: "finished client HTTP request"}
	methKV := KV{K: HTTPMethodKey, V: req.Method}
	urlKV := KV{K: HTTPURLKey, V: req.URL.String()}
	then := timeNow()
	resp, err = c.RoundTripper.RoundTrip(req)
	if err != nil {
		Error(req.Context(), err, msgKV, methKV, urlKV)
		return
	}
	ms := timeSince(then).Milliseconds()
	statusKV := KV{K: HTTPStatusKey, V: resp.Status}
	durKV := KV{K: HTTPDurationKey, V: ms}
	if c.options.iserr(resp.StatusCode) {
		if c.options.logErrBody {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				Error(req.Context(), err, msgKV, methKV, urlKV, statusKV, durKV)
				return resp, nil
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			Error(req.Context(), fmt.Errorf(resp.Status), msgKV, methKV, urlKV, statusKV, durKV, KV{K: HTTPBodyKey, V: string(body)})
		} else {
			Error(req.Context(), fmt.Errorf(resp.Status), msgKV, methKV, urlKV, statusKV, durKV)
		}
		return
	}
	Print(req.Context(), msgKV, methKV, urlKV, statusKV, durKV)
	return
}
