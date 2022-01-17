package health

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	cases := []struct {
		name           string
		deps           []Pinger
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:           "empty",
			expectedStatus: http.StatusOK,
			expectedJSON:   `{"uptime":0,"version":""}`,
		},
		{
			name:           "ok",
			deps:           singleHealthyDep("dependency"),
			expectedStatus: http.StatusOK,
			expectedJSON:   `{"uptime":0,"version":"","status":{"dependency":"OK"}}`,
		},
		{
			name:           "not ok",
			deps:           singleUnhealthyDep("dependency", fmt.Errorf("dependency is not ok")),
			expectedStatus: http.StatusServiceUnavailable,
			expectedJSON:   `{"uptime":0,"version":"","status":{"dependency":"NOT OK"}}`,
		},
		{
			name:           "multiple dependencies",
			deps:           multipleHealthyDeps("dependency1", "dependency2"),
			expectedStatus: http.StatusOK,
			expectedJSON:   `{"uptime":0,"version":"","status":{"dependency1":"OK","dependency2":"OK"}}`,
		},
		{
			name:           "multiple dependencies not ok",
			deps:           multipleUnhealthyDeps(fmt.Errorf("dependency2 is not ok"), "dependency1", "dependency2"),
			expectedStatus: http.StatusServiceUnavailable,
			expectedJSON:   `{"uptime":0,"version":"","status":{"dependency1":"OK","dependency2":"NOT OK"}}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checker := NewChecker(c.deps...)
			handler := Handler(checker)
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			if w.Code != c.expectedStatus {
				t.Errorf("got status: %d, expected %d", w.Code, c.expectedStatus)
			}
			if w.Body.String() != c.expectedJSON {
				t.Errorf("got body: %s, expected %s", w.Body.String(), c.expectedJSON)
			}
		})
	}
}
