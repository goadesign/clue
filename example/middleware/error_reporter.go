package middleware

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	goa "goa.design/goa/v3/pkg"
)

// ErrorReporter is a Goa endpoint middleware that flags a span as failed when an
// error is returned by the service method.
//
// Usage:
//
//	endpoints := gen.NewEndpoints(svc)
//	endpoints.Use(ErrorReporter())
func ErrorReporter() func(goa.Endpoint) goa.Endpoint {
	return func(e goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, req any) (any, error) {
			res, err := e(ctx, req)
			if err != nil {
				span := trace.SpanFromContext(ctx)
				span.SetStatus(codes.Error, err.Error())
			}
			return res, err
		}
	}
}
