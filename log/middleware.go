package log

import (
	"context"

	"goa.design/goa/v3/middleware"
	goa "goa.design/goa/v3/pkg"
)

// Init is a Goa endpoint middleware that initializes the logger context.
// It panics if logCtx was not initialized with Context.
//
// Usage:
//
//    ctx := log.Context(context.Background())
//    endpoints := service.NewEndpoints(svc)
//    endpoints.Use(log.Init(ctx))
func Init(logCtx context.Context) func(goa.Endpoint) goa.Endpoint {
	l := logCtx.Value(ctxLogger)
	if l == nil {
		panic("log.Init called without log.Context")
	}
	logger := l.(*logger)
	return func(e goa.Endpoint) goa.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			ctx = context.WithValue(ctx, ctxLogger, logger)
			if requestID := ctx.Value(middleware.RequestIDKey); requestID != nil {
				ctx = With(ctx, "request_id", requestID)
			}
			return e(ctx, req)
		}
	}
}
