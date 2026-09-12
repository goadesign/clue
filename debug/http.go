package debug

import (
	"net/http"

	"goa.design/clue/log"
)

// HTTP returns a middleware that manages whether debug log entries are written.
// This middleware should be used in conjunction with the MountDebugLogEnabler
// function. If the request already has a logger, it must be request-local;
// placing log.HTTP before this middleware creates one.
func HTTP() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if debugLogs.Load() {
				ctx = log.Context(ctx, log.WithDebug())
			} else {
				ctx = log.Context(ctx, log.WithNoDebug())
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
		return handler
	}
}
