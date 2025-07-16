package health

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
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
		expected := fmt.Errorf(`failed to make health check request to "dependency": Get "%s/livez": context deadline exceeded (Client.Timeout exceeded while awaiting headers)`, svr.URL)
		if err != nil && err.Error() != expected.Error() {
			t.Errorf("got: %v, expected: %v", err, expected)
		}
	})
}

func TestOptions(t *testing.T) {
	cases := []struct {
		name            string
		option          Option
		expectedScheme  string
		expectedPath    string
		expectedTimeout time.Duration
	}{
		{"default", nil, "http", "/livez", 0},
		{"scheme", WithScheme("https"), "https", "/livez", 0},
		{"path", WithPath("/healthcheck"), "http", "/healthcheck", 0},
		{"timeout", WithTimeout(10 * time.Second), "http", "/livez", 10 * time.Second},
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
			if pinger.httpClient.Timeout != c.expectedTimeout {
				t.Errorf("got timeout: %s, expected %s", pinger.httpClient.Timeout, c.expectedTimeout)
			}
		})
	}
}

func TestPingWithBaggage(t *testing.T) {
	t.Run("forwards OpenTelemetry baggage", func(t *testing.T) {
		// Set up OpenTelemetry propagator that includes baggage
		originalPropagator := otel.GetTextMapPropagator()
		defer otel.SetTextMapPropagator(originalPropagator) // restore after test

		propagator := propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
		otel.SetTextMapPropagator(propagator)

		var receivedHeaders map[string][]string
		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("got method: %s, expected GET", r.Method)
			}
			if r.URL.Path != "/livez" {
				t.Errorf("got path: %s, expected /livez", r.URL.Path)
			}
			// Capture headers to verify baggage was forwarded
			receivedHeaders = r.Header
			w.WriteHeader(http.StatusOK)
		}
		svr := httptest.NewServer(http.HandlerFunc(handler))
		defer svr.Close()
		u, _ := url.Parse(svr.URL)
		pinger := NewPinger("dependency", u.Host)

		// Create baggage with test data
		property, err := baggage.NewKeyValuePropertyRaw("property", "test-value")
		if err != nil {
			t.Fatalf("failed to create baggage property: %v", err)
		}
		member, err := baggage.NewMemberRaw("test-key", "test-member-value", property)
		if err != nil {
			t.Fatalf("failed to create baggage member: %v", err)
		}
		bag, err := baggage.New(member)
		if err != nil {
			t.Fatalf("failed to create baggage: %v", err)
		}

		// Add baggage to context
		ctx := baggage.ContextWithBaggage(context.Background(), bag)

		// Make the ping request
		err = pinger.Ping(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verify baggage header was sent
		if receivedHeaders == nil {
			t.Fatal("no headers received")
		}

		baggageHeaders, exists := receivedHeaders["Baggage"]
		if !exists {
			t.Error("baggage header not found in request")
		} else if len(baggageHeaders) == 0 {
			t.Error("baggage header is empty")
		} else {
			// Verify baggage content using propagation to extract and validate
			extractedCtx := propagator.Extract(context.Background(), propagation.HeaderCarrier(receivedHeaders))
			extractedBag := baggage.FromContext(extractedCtx)

			testMember := extractedBag.Member("test-key")
			if testMember.Value() != "test-member-value" {
				t.Errorf("got baggage member value: %s, expected test-member-value", testMember.Value())
			}

			properties := testMember.Properties()
			if len(properties) != 1 {
				t.Errorf("got %d properties, expected 1", len(properties))
			} else {
				prop := properties[0]
				if prop.Key() != "property" {
					t.Errorf("got property key: %s, expected property", prop.Key())
				}
				if val, ok := prop.Value(); !ok || val != "test-value" {
					t.Errorf("got property value: %s (ok=%t), expected test-value", val, ok)
				}
			}
		}
	})

	t.Run("works without baggage", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}
		svr := httptest.NewServer(http.HandlerFunc(handler))
		defer svr.Close()
		u, _ := url.Parse(svr.URL)
		pinger := NewPinger("dependency", u.Host)

		// Make ping request without baggage
		err := pinger.Ping(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
