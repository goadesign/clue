package log

import (
	"context"
	"net/http"
	"regexp"

	"goa.design/goa/v3/middleware"
)

type (
	// HTTPOption is a function that applies a configuration option to log
	// HTTP middleware.
	HTTPOption func(*httpOptions)

	httpOptions struct {
		pathFilters []*regexp.Regexp
	}
)

// HTTP returns a HTTP middleware that initializes the logger context.  It
// panics if logCtx was not initialized with Context.
func HTTP(logCtx context.Context, opts ...HTTPOption) func(http.Handler) http.Handler {
	MustContainLogger(logCtx)
	var options httpOptions
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

// WithPathFilter adds a path filter to the HTTP middleware. Requests whose path
// match the filter are not logged. WithPathFilter can be called multiple times
// to add multiple filters.
func WithPathFilter(filter *regexp.Regexp) HTTPOption {
	return func(o *httpOptions) {
		o.pathFilters = append(o.pathFilters, filter)
	}
}
