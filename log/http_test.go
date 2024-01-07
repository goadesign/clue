package log

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	goa "goa.design/goa/v3/pkg"
)

func TestHTTP(t *testing.T) {
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()
	timeSince = func(_ time.Time) time.Duration { return 42 * time.Millisecond }
	defer func() { timeSince = time.Since }()
	shortID = func() string { return "test-request-id" }
	defer func() { shortID = randShortID }()

	prefix := `{"time":"2022-01-09T20:29:45Z","level":"info","request_id":"test-request-id","msg":"start","http.method":"GET","http.url":"http://example.com","http.remote_addr":""}`
	entry := `{"time":"2022-01-09T20:29:45Z","level":"info","request_id":"test-request-id","key1":"value1","key2":"value2"}`
	suffix := `{"time":"2022-01-09T20:29:45Z","level":"info","request_id":"test-request-id","msg":"end","http.method":"GET","http.url":"http://example.com","http.status":0,"http.time_ms":42,"http.bytes":0}`

	cases := []struct {
		name     string
		opt      HTTPLogOption
		expected string
	}{
		{
			name:     "default",
			expected: prefix + "\n" + entry + "\n" + suffix + "\n",
		},
		{
			name: "with path filter",
			opt:  WithPathFilter(regexp.MustCompile("")),
		},
		{
			name:     "with disable request logging",
			opt:      WithDisableRequestLogging(),
			expected: entry + "\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var handler http.Handler = http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
				Print(req.Context(), KV{"key1", "value1"}, KV{"key2", "value2"})
			})
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatJSON))

			handler = HTTP(ctx, c.opt)(handler)

			req, _ := http.NewRequest("GET", "http://example.com", nil)

			handler.ServeHTTP(nil, req)

			assert.Equal(t, c.expected, buf.String())
		})
	}
}

func TestEndpoint(t *testing.T) {
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()
	endpoint := func(ctx context.Context, _ interface{}) (interface{}, error) {
		Printf(ctx, "test")
		return nil, nil
	}
	cases := []struct {
		name     string
		sname    string
		mname    string
		expected string
	}{
		{"service and method name", "Service", "Method", `{"time":"2022-01-09T20:29:45Z","level":"info","goa.service":"Service","goa.method":"Method","msg":"test"}`},
		{"no service name", "", "Method", `{"time":"2022-01-09T20:29:45Z","level":"info","goa.method":"Method","msg":"test"}`},
		{"no method name", "Service", "", `{"time":"2022-01-09T20:29:45Z","level":"info","goa.service":"Service","msg":"test"}`},
		{"no service or method name", "", "", `{"time":"2022-01-09T20:29:45Z","level":"info","msg":"test"}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatJSON))
			if c.sname != "" {
				ctx = context.WithValue(ctx, goa.ServiceKey, c.sname)
			}
			if c.mname != "" {
				ctx = context.WithValue(ctx, goa.MethodKey, c.mname)
			}

			_, err := Endpoint(endpoint)(ctx, nil)
			assert.NoError(t, err)

			assert.Equal(t, c.expected+"\n", buf.String())
		})
	}
}

func TestClient(t *testing.T) {
	successLogs := `time=2022-01-09T20:29:45Z level=info msg="finished client HTTP request" http.method=GET http.url=$URL http.status="200 OK" http.time_ms=42`
	errorLogs := `time=2022-01-09T20:29:45Z level=error err=error msg="finished client HTTP request" http.method=GET http.url=$URL`
	statusLogs := `time=2022-01-09T20:29:45Z level=error err="200 OK" msg="finished client HTTP request" http.method=GET http.url=$URL http.status="200 OK" http.time_ms=42`
	cases := []struct {
		name      string
		noLog     bool
		clientErr error
		opt       HTTPClientLogOption
		expected  string
	}{
		{"no logger", true, nil, nil, ""},
		{"success", false, nil, nil, successLogs},
		{"error", false, fmt.Errorf("error"), nil, errorLogs},
		{"error with status", false, nil, WithErrorStatus(200), statusLogs},
	}
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()
	duration := 42 * time.Millisecond
	since := timeSince
	timeSince = func(_ time.Time) time.Duration { return duration }
	defer func() { timeSince = since }()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf))
			if c.noLog {
				ctx = context.Background()
			}
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) { rw.Write([]byte(`OK`)) })) //nolint:errcheck
			defer server.Close()
			client := server.Client()
			if c.clientErr != nil {
				client.Transport = &errorClient{err: c.clientErr}
			}
			if c.opt != nil {
				client.Transport = Client(client.Transport, c.opt)
			} else {
				client.Transport = Client(client.Transport)
			}

			req, _ := http.NewRequest("GET", server.URL, nil)
			req = req.WithContext(ctx)

			client.Do(req) // nolint:errcheck

			expected := strings.ReplaceAll(c.expected, "$URL", server.URL)
			assert.Equal(t, strings.TrimSpace(buf.String()), expected)
		})
	}
}

func TestWithPathFilter(t *testing.T) {
	var handler http.Handler = http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		Print(req.Context(), KV{"key1", "value1"}, KV{"key2", "value2"})
	})
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatJSON))

	handler = HTTP(ctx, WithPathFilter(regexp.MustCompile("/path/to/ignore")))(handler)

	req, _ := http.NewRequest("GET", "http://example.com/path/to/ignore", nil)

	handler.ServeHTTP(nil, req)

	assert.Empty(t, buf.String())
}

type errorClient struct {
	err error
}

func (c *errorClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, c.err
}
