package debug

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goa.design/clue/internal/testsvc"
	"goa.design/clue/internal/testsvc/gen/test"
	"goa.design/clue/log"
	goahttp "goa.design/goa/v3/http"
)

func TestMountDebugLogEnabler(t *testing.T) {
	cases := []struct {
		name         string
		prefix       string
		query        string
		onval        string
		offval       string
		url          string
		expectedResp string
	}{
		{"defaults", "", "", "", "", "/debug", `{"debug-logs":"off"}`},
		{"defaults-enable", "", "", "", "", "/debug?debug-logs=on", `{"debug-logs":"on"}`},
		{"defaults-disable", "", "", "", "", "/debug?debug-logs=off", `{"debug-logs":"off"}`},
		{"prefix", "test", "", "", "", "/test", `{"debug-logs":"off"}`},
		{"prefix-enable", "test", "", "", "", "/test?debug-logs=on", `{"debug-logs":"on"}`},
		{"prefix-disable", "test", "", "", "", "/test?debug-logs=off", `{"debug-logs":"off"}`},
		{"query", "", "debug", "", "", "/debug", `{"debug":"off"}`},
		{"query-enable", "", "debug", "", "", "/debug?debug=on", `{"debug":"on"}`},
		{"query-disable", "", "debug", "", "", "/debug?debug=off", `{"debug":"off"}`},
		{"onval-enable", "", "", "foo", "", "/debug?debug-logs=foo", `{"debug-logs":"foo"}`},
		{"offval-disable", "", "", "", "bar", "/debug?debug-logs=bar", `{"debug-logs":"bar"}`},
		{"prefix-query-enable", "test", "debug", "", "", "/test?debug=on", `{"debug":"on"}`},
		{"prefix-query-disable", "test", "debug", "", "", "/test?debug=off", `{"debug":"off"}`},
		{"prefix-onval-enable", "test", "", "foo", "", "/test?debug-logs=foo", `{"debug-logs":"foo"}`},
		{"prefix-offval-disable", "test", "", "", "bar", "/test?debug-logs=bar", `{"debug-logs":"bar"}`},
		{"prefix-query-onval-enable", "test", "debug", "foo", "", "/test?debug=foo", `{"debug":"foo"}`},
		{"prefix-query-offval-disable", "test", "debug", "", "bar", "/test?debug=bar", `{"debug":"bar"}`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mux := http.NewServeMux()
			var options []DebugLogEnablerOption
			if c.prefix != "" {
				options = append(options, WithPath(c.prefix))
			}
			if c.query != "" {
				options = append(options, WithQuery(c.query))
			}
			if c.onval != "" {
				options = append(options, WithOnValue(c.onval))
			}
			if c.offval != "" {
				options = append(options, WithOffValue(c.offval))
			}
			MountDebugLogEnabler(mux, options...)
			ts := httptest.NewServer(mux)
			defer ts.Close()

			status, resp := makeRequest(t, ts.URL+c.url)

			assert.Equal(t, http.StatusOK, status)
			assert.Equal(t, c.expectedResp, resp)
		})
	}
}

func TestMountPprofHandlers(t *testing.T) {
	mux := goahttp.NewMuxer()
	MountPprofHandlers(Adapt(mux))
	MountPprofHandlers(Adapt(mux), WithPrefix("test"))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK")) // nolint: errcheck
	})
	mux.Handle(http.MethodGet, "/", handler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	status, resp := makeRequest(t, ts.URL)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "OK", resp)

	paths := []string{
		"/debug/pprof/",
		"/test/",
		"/debug/pprof/allocs",
		"/debug/pprof/block",
		"/debug/pprof/cmdline",
		"/debug/pprof/goroutine",
		"/debug/pprof/heap",
		"/debug/pprof/mutex",
		// "/debug/pprof/profile?seconds=1", # Takes too long to run on each test
		"/debug/pprof/symbol",
		"/debug/pprof/threadcreate",
		// "/debug/pprof/trace", # Takes too long to run on each test
	}
	for _, path := range paths {
		status, resp = makeRequest(t, ts.URL+path)
		assert.Equal(t, http.StatusOK, status)
		assert.NotEmpty(t, resp)
	}
}

func TestDebugPayloads(t *testing.T) {
	svc := testsvc.Service{}
	var buf bytes.Buffer
	vals, vali := "test", 1
	payload := &test.Fields{S: &vals, I: &vali}
	testErr := errors.New("test error")
	newLogContext := func() context.Context {
		return log.Context(context.Background(), log.WithOutput(&buf), log.WithFormat(logKeyValsOnly))
	}
	newDebugLogContext := func() context.Context {
		return log.Context(newLogContext(), log.WithDebug())
	}
	formatTest := func(_ context.Context, a interface{}) string {
		return "test"
	}

	cases := []struct {
		name         string
		ctx          context.Context
		option       LogPayloadsOption
		methErr      error
		expectedLogs string
		expectedErr  string
	}{
		{"no debug", newLogContext(), nil, nil, "", ""},
		{"debug", newDebugLogContext(), nil, nil, `payload={"S":"test","I":1} result={"S":"test","I":1} `, ""},
		{"error", newLogContext(), nil, testErr, "", "test error"},
		{"debug error", newDebugLogContext(), nil, testErr, `payload={"S":"test","I":1} `, "test error"},
		{"maxsize", newDebugLogContext(), WithMaxSize(1), nil, `payload={ result={ `, ""},
		{"format", newDebugLogContext(), WithFormat(formatTest), nil, `payload=test result=test `, ""},
		{"client", newDebugLogContext(), WithClient(), nil, `client-payload={"S":"test","I":1} client-result={"S":"test","I":1} `, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf.Reset()
			svc.HTTPFunc = func(_ context.Context, f *testsvc.Fields) (*testsvc.Fields, error) {
				if c.expectedErr != "" {
					return nil, errors.New(c.expectedErr)
				}
				return f, nil
			}
			endpoint := test.NewHTTPMethodEndpoint(&svc)
			endpoint = LogPayloads(c.option)(endpoint)
			res, err := endpoint(c.ctx, payload)
			assert.Equal(t, c.expectedLogs, buf.String())
			if err != nil {
				assert.Equal(t, c.expectedErr, err.Error())
				return
			}
			require.Empty(t, c.expectedErr)
			require.NotNil(t, res)
			assert.Equal(t, payload, res)
		})
	}
}

func logKeyValsOnly(e *log.Entry) []byte {
	var buf bytes.Buffer
	for _, kv := range e.KeyVals {
		buf.WriteString(kv.K)
		buf.WriteString("=")
		buf.WriteString(fmt.Sprintf("%v", kv.V))
		buf.WriteString(" ")
	}
	return buf.Bytes()
}

func makeRequest(t *testing.T, url string) (int, string) {
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make request %q: %s", url, err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			t.Logf("failed to close response body: %v", err)
		}
	}()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body) // nolint: errcheck
	return res.StatusCode, buf.String()
}
