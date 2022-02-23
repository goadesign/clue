package health

import (
	"encoding/json"
	"net/http"
)

// Handler returns a HTTP handler that serves health check requests. The
// response body is the JSON encoded health status returned by chk.Check(). The
// response status is 200 if chk.Check() returns a nil error, 503 otherwise.
func Handler(chk Checker) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, healthy := chk.Check(r.Context())
		b, _ := json.Marshal(h)
		if healthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		w.Write(b)
	})
}
