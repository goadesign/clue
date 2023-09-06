package debug

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"goa.design/clue/log"
)

func TestHTTP(t *testing.T) {
	// Create log context
	var buf bytes.Buffer
	ctx := log.Context(context.Background(), log.WithOutput(&buf), log.WithFormat(logKeyValsOnly))
	log.FlushAndDisableBuffering(ctx)

	// Create HTTP handler
	mux := http.NewServeMux()
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Info(r.Context(), log.KV{K: "test", V: "info"})
		log.Debug(r.Context(), log.KV{K: "test", V: "debug"})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK")) // nolint: errcheck
	})

	// Mount debug handler and log middleware
	MountDebugLogEnabler(mux)
	handler = HTTP()(handler)
	handler = log.HTTP(ctx, log.WithDisableRequestLogging())(handler)

	// Start test server
	mux.Handle("/", handler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	steps := []struct {
		name         string
		on           bool
		off          bool
		expectedResp string
		expectedLogs string
	}{
		{"start", false, false, "", "test=info "},
		{"turn debug logs on", true, false, `{"debug-logs":true}`, "test=info test=debug "},
		{"with debug logs on", false, false, `{"debug-logs":true}`, "test=info test=debug "},
		{"turn debug logs off", false, true, `{"debug-logs":false}`, "test=info "},
		{"with debug logs off", false, false, `{"debug-logs":false}`, "test=info "},
	}
	for _, step := range steps {
		if step.on {
			makeRequest(t, ts.URL+"/debug?debug-logs=on")
		}
		if step.off {
			makeRequest(t, ts.URL+"/debug?debug-logs=off")
		}

		status, resp := makeRequest(t, ts.URL)

		if status != http.StatusOK {
			t.Errorf("%s: got status %d, expected %d", step.name, status, http.StatusOK)
		}
		if resp != "OK" {
			t.Errorf("%s: got body %q, expected %q", step.name, resp, "OK")
		}
		if buf.String() != step.expectedLogs {
			t.Errorf("%s: got logs %q, expected %q", step.name, buf.String(), step.expectedLogs)
		}
		buf.Reset()
	}
}
