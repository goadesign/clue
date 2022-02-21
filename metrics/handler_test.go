package metrics

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"goa.design/clue/log"
)

func TestHandler(t *testing.T) {
	errTest := errors.New("test")
	cases := []struct {
		name             string
		err              error
		expectedLog      string
		expectedResponse string
	}{
		{"no error", nil, "", "test_metric 1"},
		{"error", errTest, errTest.Error(), ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := log.Context(context.Background(),
				log.WithOutput(&buf),
				log.WithFormat(func(e *log.Entry) []byte { return []byte(e.KeyVals[0].V.(string)) }))
			reg := NewTestRegistry(t)
			gat := &mockGatherer{err: c.err}
			req, _ := http.NewRequest("GET", "/metrics", nil)
			w := httptest.NewRecorder()

			Handler(ctx, WithGatherer(gat), WithHandlerRegisterer(reg)).ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if c.err != nil {
				if !strings.Contains(string(body), c.err.Error()) {
					t.Errorf("expected error message %q, got %q", c.err.Error(), string(body))
				}
				if w.Code != http.StatusInternalServerError {
					t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
				}
				if !strings.Contains(buf.String(), c.expectedLog) {
					t.Errorf("expected log %q, got %q", c.expectedLog, buf.String())
				}
				return
			}
			if buf.Len() > 0 {
				t.Errorf("unexpected log message %q", buf.String())
			}
			if !strings.Contains(string(body), c.expectedResponse) {
				t.Errorf("expected response %q, got %q", c.expectedResponse, string(body))
			}
		})
	}
}

type mockGatherer struct {
	err error
}

func (m *mockGatherer) Gather() ([]*dto.MetricFamily, error) {
	var (
		name         = "test_metric"
		one  float64 = 1
		typ          = dto.MetricType_COUNTER
	)
	if m.err != nil {
		return nil, m.err
	}
	return []*dto.MetricFamily{{
		Name:   &name,
		Type:   &typ,
		Metric: []*dto.Metric{{Counter: &dto.Counter{Value: &one}}},
	}}, nil
}
