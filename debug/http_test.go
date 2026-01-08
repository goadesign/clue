package debug

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"goa.design/clue/log"
)

func TestHTTP(t *testing.T) {
	// Create log context
	var buf bytes.Buffer
	ctx := log.Context(context.Background(),
		log.WithOutputs(log.Output{Writer: &buf, Format: logKeyValsOnly}),
	)
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
	handler = log.HTTP(ctx,
		log.WithDisableRequestLogging(),
		log.WithDisableRequestID())(handler)

	// Start test server
	mux.Handle("/", handler)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	steps := []struct {
		name    string
		on      bool
		off     bool
		wantLog string
	}{
		{"start", false, false, "test=info "},
		{"turn debug logs on", true, false, "test=info test=debug "},
		{"with debug logs on", false, false, "test=info test=debug "},
		{"turn debug logs off", false, true, "test=info "},
		{"with debug logs off", false, false, "test=info "},
	}
	for _, step := range steps {
		if step.on {
			makeRequest(t, ts.URL+"/debug?debug-logs=on")
		}
		if step.off {
			makeRequest(t, ts.URL+"/debug?debug-logs=off")
		}

		status, resp := makeRequest(t, ts.URL)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "OK", resp)
		assert.Equal(t, step.wantLog, buf.String())
		buf.Reset()
	}
}
