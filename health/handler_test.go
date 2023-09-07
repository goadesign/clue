package health

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	type (
		assertFunc func(*testing.T, string)
	)
	var (
		assertContains = func(contains string) assertFunc {
			return func(t *testing.T, body string) {
				assert.Contains(t, body, contains)
			}
		}
		assertEqual = func(expected string) assertFunc {
			return func(t *testing.T, body string) {
				assert.Equal(t, expected, body)
			}
		}
		assertNotContains = func(contains string) assertFunc {
			return func(t *testing.T, body string) {
				assert.NotContains(t, body, contains)
			}
		}
	)
	cases := []struct {
		name                 string
		deps                 []Pinger
		expectedStatus       int
		assertFuncsPerAccept map[string][]assertFunc
	}{
		{
			name:           "empty",
			expectedStatus: http.StatusOK,
			assertFuncsPerAccept: map[string][]assertFunc{
				"":                {assertEqual(`{"uptime":0,"version":""}`)},
				"application/xml": {assertEqual("<health><uptime>0</uptime><version></version></health>")},
				"application/gob": {assertNotContains("status"), assertNotContains("dependency")},
			},
		},
		{
			name:           "ok",
			deps:           singleHealthyDep("dependency"),
			expectedStatus: http.StatusOK,
			assertFuncsPerAccept: map[string][]assertFunc{
				"":                {assertEqual(`{"uptime":0,"version":"","status":{"dependency":"OK"}}`)},
				"application/xml": {assertEqual("<health><uptime>0</uptime><version></version><status><dependency>OK</dependency></status></health>")},
				"application/gob": {assertContains("Status"), assertContains("dependency"), assertNotContains("NOT OK")},
			},
		},
		{
			name:           "not ok",
			deps:           singleUnhealthyDep("dependency", fmt.Errorf("dependency is not ok")),
			expectedStatus: http.StatusServiceUnavailable,
			assertFuncsPerAccept: map[string][]assertFunc{
				"":                {assertEqual(`{"uptime":0,"version":"","status":{"dependency":"NOT OK"}}`)},
				"application/xml": {assertEqual("<health><uptime>0</uptime><version></version><status><dependency>NOT OK</dependency></status></health>")},
				"application/gob": {assertContains("Status"), assertContains("dependency"), assertContains("NOT OK")},
			},
		},
		{
			name:           "multiple dependencies",
			deps:           multipleHealthyDeps("dependency1", "dependency2"),
			expectedStatus: http.StatusOK,
			assertFuncsPerAccept: map[string][]assertFunc{
				"":                {assertEqual(`{"uptime":0,"version":"","status":{"dependency1":"OK","dependency2":"OK"}}`)},
				"application/xml": {assertEqual("<health><uptime>0</uptime><version></version><status><dependency1>OK</dependency1><dependency2>OK</dependency2></status></health>")},
				"application/gob": {assertContains("dependency1"), assertContains("dependency2"), assertNotContains("NOT OK")},
			},
		},
		{
			name:           "multiple dependencies not ok",
			deps:           multipleUnhealthyDeps(fmt.Errorf("dependency2 is not ok"), "dependency1", "dependency2"),
			expectedStatus: http.StatusServiceUnavailable,
			assertFuncsPerAccept: map[string][]assertFunc{
				"":                {assertEqual(`{"uptime":0,"version":"","status":{"dependency1":"OK","dependency2":"NOT OK"}}`)},
				"application/xml": {assertEqual("<health><uptime>0</uptime><version></version><status><dependency1>OK</dependency1><dependency2>NOT OK</dependency2></status></health>")},
				"application/gob": {assertContains("dependency1"), assertContains("dependency2"), assertContains("NOT OK")},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checker := NewChecker(c.deps...)
			handler := Handler(checker)
			req := httptest.NewRequest("GET", "/", nil)
			for accept, fns := range c.assertFuncsPerAccept {
				t.Run("Accept:"+accept, func(t *testing.T) {
					req.Header.Set("Accept", accept)
					w := httptest.NewRecorder()
					handler.ServeHTTP(w, req)
					if w.Code != c.expectedStatus {
						t.Errorf("got status: %d, expected %d", w.Code, c.expectedStatus)
					}
					body := strings.TrimSpace(w.Body.String())
					for _, fn := range fns {
						fn(t, body)
					}
				})
			}
		})
	}
}
