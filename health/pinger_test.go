package health

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	cases := []struct {
		name     string
		status   int
		expected error
	}{
		{
			name:     "ok",
			status:   http.StatusOK,
			expected: nil,
		},
		{
			name:     "not ok",
			status:   http.StatusServiceUnavailable,
			expected: fmt.Errorf(`health-check for "dependency" returned status 503`),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("got method: %s, expected GET", r.Method)
				}
				if r.URL.Path != "/livez" {
					t.Errorf("got path: %s, expected /livez", r.URL.Path)
				}
				w.WriteHeader(c.status)
			}
			svr := httptest.NewServer(http.HandlerFunc(handler))
			defer svr.Close()
			u, _ := url.Parse(svr.URL)
			pinger := NewPinger("dependency", u.Host)
			if pinger.Name() != "dependency" {
				t.Errorf("got name: %s, expected dependency", pinger.Name())
			}
			err := pinger.Ping(context.Background())
			if err == nil && c.expected != nil {
				t.Errorf("got error: %v, expected %v", err, c.expected)
				return
			}
			if err != nil && c.expected == nil {
				t.Errorf("got error: %v, expected %v", err, c.expected)
				return
			}
			if err != nil && err.Error() != c.expected.Error() {
				t.Errorf("got: %v, expected: %v", err, c.expected)
			}
		})
	}
}

func TestOptions(t *testing.T) {
	cases := []struct {
		name           string
		option         Option
		expectedScheme string
		expectedPath   string
	}{
		{"default", nil, "http", "/livez"},
		{"scheme", WithScheme("https"), "https", "/livez"},
		{"path", WithPath("/healthcheck"), "http", "/healthcheck"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var pinger *client
			if c.option == nil {
				pinger = NewPinger("dependency", "localhost").(*client)
			} else {
				pinger = NewPinger("dependency", "localhost", c.option).(*client)
			}
			if pinger.Name() != "dependency" {
				t.Errorf("got name: %s, expected dependency", pinger.Name())
			}
			if pinger.req.URL.Scheme != c.expectedScheme {
				t.Errorf("got scheme: %s, expected %s", pinger.req.URL.Scheme, c.expectedScheme)
			}
			if pinger.req.URL.Path != c.expectedPath {
				t.Errorf("got path: %s, expected %s", pinger.req.URL.Path, c.expectedPath)
			}
		})
	}
}
