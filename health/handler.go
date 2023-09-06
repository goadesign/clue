package health

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
)

// Handler returns a HTTP handler that serves health check requests. The
// response body is the JSON encoded health status returned by chk.Check(). The
// response status is 200 if chk.Check() returns a nil error, 503 otherwise.
func Handler(chk Checker) http.HandlerFunc {
	encoder := goahttp.ResponseEncoder
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		enc := encoder(ctx, w)
		h, healthy := chk.Check(ctx)
		if healthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		enc.Encode(h) // nolint: errcheck
	})
}
