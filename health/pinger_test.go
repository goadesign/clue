package health

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	t.Run("timeout", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}
		svr := httptest.NewServer(http.HandlerFunc(handler))
		defer svr.Close()
		u, _ := url.Parse(svr.URL)
		pinger := NewPinger("dependency", u.Host, WithTimeout(time.Second))
		if pinger.Name() != "dependency" {
			t.Errorf("got name: %s, expected dependency", pinger.Name())
		}
		err := pinger.Ping(context.Background())
		// The full error message can vary due to some race conditions on child context cancellation. Check the prefix.
		expectedPrefix := fmt.Sprintf(`failed to make health check request to "dependency": Get "%s/livez":`, svr.URL)
		if err != nil && !strings.Contains(err.Error(), expectedPrefix) {
			t.Errorf("got: %v, expected prefix of: %v", err, expectedPrefix)
		}
	})
}

func TestOptions(t *testing.T) {
	cases := []struct {
		name              string
		option            Option
		expectedScheme    string
		expectedPath      string
		expectedTimeout   time.Duration
		expectedTransport http.RoundTripper
	}{
		{"default", nil, "http", "/livez", 0, nil},
		{"scheme", WithScheme("https"), "https", "/livez", 0, nil},
		{"path", WithPath("/healthcheck"), "http", "/healthcheck", 0, nil},
		{"timeout", WithTimeout(10 * time.Second), "http", "/livez", 10 * time.Second, nil},
		{"transport", WithTransport(http.DefaultTransport), "http", "/livez", 0, http.DefaultTransport},
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
			u, err := url.Parse(pinger.httpURL)
			if err != nil {
				t.Errorf("got error: %v, expected %v", err, nil)
			}
			if u.Scheme != c.expectedScheme {
				t.Errorf("got scheme: %s, expected %s", u.Scheme, c.expectedScheme)
			}
			if u.Path != c.expectedPath {
				t.Errorf("got path: %s, expected %s", u.Path, c.expectedPath)
			}
			if pinger.httpClient.Timeout != c.expectedTimeout {
				t.Errorf("got timeout: %s, expected %s", pinger.httpClient.Timeout, c.expectedTimeout)
			}
			if c.expectedTransport != nil && pinger.httpClient.Transport != c.expectedTransport {
				t.Errorf("got transport: %v, expected %v", pinger.httpClient.Transport, c.expectedTransport)
			}
		})
	}
}

func TestOptions_DefaultTransport(t *testing.T) {
	pinger := NewPinger("dependency", "localhost").(*client)
	if _, ok := pinger.httpClient.Transport.(*otelhttp.Transport); !ok {
		t.Errorf("got transport: %v, expected otelhttp.Transport", pinger.httpClient.Transport)
	}
}
