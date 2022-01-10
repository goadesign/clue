package log

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// SetContext is a Goa endpoint middleware that initializes the logger context.
//
// Usage:
//
//    endpoints := service.NewEndpoints(svc)
//    endpoints.Use(log.SetContext(log.WithFormat(log.FormatJSON)))
func SetContext(opts ...LogOption) func(goa.Endpoint) goa.Endpoint {
	return func(e goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			ctx = Context(ctx, opts...)
			return e(ctx, req)
		}
	}
}
