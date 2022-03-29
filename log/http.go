package log

import (
	"context"
	"net/http"

	"goa.design/goa/v3/middleware"
)

// HTTP returns a HTTP middleware that initializes the logger context.  It
// panics if logCtx was not initialized with Context.
func HTTP(logCtx context.Context) func(http.Handler) http.Handler {
	MustContainLogger(logCtx)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := WithContext(req.Context(), logCtx)
			if requestID := req.Context().Value(middleware.RequestIDKey); requestID != nil {
				ctx = With(ctx, KV{"requestID", requestID})
			}
			h.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
