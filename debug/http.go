package debug

import (
	"net/http"

	"goa.design/clue/log"
)

// HTTP returns a middleware that manages whether debug log entries are written.
// This middleware should be used in conjunction with the MountDebugLogEnabler
// function.
func HTTP() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if debugLogs {
				ctx := log.Context(r.Context(), log.WithDebug())
				r = r.WithContext(ctx)
			} else {
				ctx := log.Context(r.Context(), log.WithNoDebug())
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
		return handler
	}
}
