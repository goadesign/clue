package debug

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"goa.design/clue/internal/testsvc"
	"goa.design/clue/internal/testsvc/gen/test"
	"goa.design/clue/log"
)

func TestMountDebugLogEnabler(t *testing.T) {
	mux := http.NewServeMux()
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Info(r.Context(), log.KV{K: "test", V: "info"})
		log.Debug(r.Context(), log.KV{K: "test", V: "debug"})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	handler = MountDebugLogEnabler("/debug", mux)(handler)
	var buf bytes.Buffer
	ctx := log.Context(context.Background(), log.WithOutput(&buf), log.WithFormat(logKeyValsOnly))
	log.FlushAndDisableBuffering(ctx)
	handler = log.HTTP(ctx)(handler)
	mux.Handle("/", handler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	steps := []struct {
		name         string
		enable       bool
		disable      bool
		expectedResp string
		expectedLogs string
	}{
		{"default", false, false, "", "test=info "},
		{"enable debug", true, false, `{"debug-logs":true}`, "debug-logs=true test=info test=debug "},
		{"disable debug", false, true, `{"debug-logs":false}`, "debug-logs=false test=info "},
	}

	for _, c := range steps {
		if c.enable {
			status, resp := makeRequest(t, ts.URL+"/debug?enable=on")
			if status != http.StatusOK {
				t.Errorf("%s: got status %d, expected %d", c.name, status, http.StatusOK)
			}
			if resp != c.expectedResp {
				t.Errorf("%s: got body %q, expected %q", c.name, resp, c.expectedResp)
			}
		}
		if c.disable {
			status, resp := makeRequest(t, ts.URL+"/debug?enable=off")
			if status != http.StatusOK {
				t.Errorf("%s: got status %d, expected %d", c.name, status, http.StatusOK)
			}
			if resp != c.expectedResp {
				t.Errorf("%s: got body %q, expected %q", c.name, resp, c.expectedResp)
			}
		}
		buf.Reset()

		status, resp := makeRequest(t, ts.URL)
		if status != http.StatusOK {
			t.Errorf("%s: got status %d, expected %d", c.name, status, http.StatusOK)
		}
		if resp != "OK" {
			t.Errorf("%s: got body %q, expected %q", c.name, resp, "OK")
		}
		if buf.String() != c.expectedLogs {
			t.Errorf("%s: got logs %q, expected %q", c.name, buf.String(), c.expectedLogs)
		}
	}
}

func TestMountPprofHandlers(t *testing.T) {
	mux := http.NewServeMux()
	MountPprofHandlers(mux)
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.Handle("/", handler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	status, resp := makeRequest(t, ts.URL)
	if status != http.StatusOK {
		t.Errorf("got status %d, expected %d", status, http.StatusOK)
	}
	if resp != "OK" {
		t.Errorf("got body %q, expected %q", resp, "OK")
	}

	paths := []string{
		"/debug/pprof/",
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
		if status != http.StatusOK {
			t.Errorf("got status %d, expected %d", status, http.StatusOK)
		}
		if resp == "" {
			t.Errorf("got body %q, expected non-empty", resp)
		}
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
			svc.HTTPFunc = func(ctx context.Context, f *testsvc.Fields) (*testsvc.Fields, error) {
				if c.expectedErr != "" {
					return nil, errors.New(c.expectedErr)
				}
				return f, nil
			}
			endpoint := test.NewHTTPMethodEndpoint(&svc)
			endpoint = LogPayloads(c.option)(endpoint)
			res, err := endpoint(c.ctx, payload)
			if buf.String() != c.expectedLogs {
				t.Errorf("got unexpected logs %q", buf.String())
			}
			if err != nil {
				if err.Error() != c.expectedErr {
					t.Errorf("got unexpected error %v", err)
				}
				return
			}
			if c.expectedErr != "" {
				t.Fatalf("expected error %q", c.expectedErr)
			}
			if res == nil {
				t.Fatal("got nil response")
			}
			if *(res.(*test.Fields)) != *payload {
				t.Errorf("got unexpected response %v", res)
			}
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
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return res.StatusCode, buf.String()
}
