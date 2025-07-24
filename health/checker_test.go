package health

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestCheck(t *testing.T) {
	cases := []struct {
		name            string
		deps            []Pinger
		expectedStatus  map[string]string
		expectedHealthy bool
	}{
		{
			name:            "empty",
			expectedStatus:  map[string]string{},
			expectedHealthy: true,
		},
		{
			name:            "ok",
			deps:            singleHealthyDep("dependency"),
			expectedStatus:  map[string]string{"dependency": "OK"},
			expectedHealthy: true,
		},
		{
			name:            "not ok",
			deps:            singleUnhealthyDep("dependency", fmt.Errorf("dependency is not ok")),
			expectedStatus:  map[string]string{"dependency": "NOT OK"},
			expectedHealthy: false,
		},
		{
			name: "multiple dependencies",
			deps: multipleHealthyDeps("dependency1", "dependency2"),
			expectedStatus: map[string]string{
				"dependency1": "OK",
				"dependency2": "OK",
			},
			expectedHealthy: true,
		},
		{
			name: "multiple dependencies not ok",
			deps: multipleUnhealthyDeps(fmt.Errorf("dependency2 is not ok"), "dependency1", "dependency2"),
			expectedStatus: map[string]string{
				"dependency1": "OK",
				"dependency2": "NOT OK",
			},
			expectedHealthy: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checker := NewChecker(c.deps...)
			res, err := checker.Check(context.Background())
			if err != c.expectedHealthy {
				t.Errorf("unexpected error: %v", err)
			}
			if res.Uptime != int64(time.Since(StartedAt).Seconds()) {
				t.Errorf("unexpected uptime: %d", res.Uptime)
			}
			if res.Version != Version {
				t.Errorf("unexpected version: %s", res.Version)
			}
			if len(res.Status) != len(c.expectedStatus) {
				t.Errorf("unexpected status: %v", res.Status)
			}
			for k, v := range c.expectedStatus {
				if res.Status[k] != v {
					t.Errorf("unexpected status for %s: %s", k, res.Status[k])
				}
			}
		})
	}
}

func TestCheck_Tracing(t *testing.T) {
	// Set up test tracer
	exporter := tracetest.NewInMemoryExporter()
	tp := trace.NewTracerProvider(trace.WithSyncer(exporter))
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	// Create a parent span
	tracer := tp.Tracer("test")
	parentCtx, parentSpan := tracer.Start(context.Background(), "parent-span")
	defer parentSpan.End()

	// Create checker with multiple dependencies
	deps := []Pinger{
		&mockDep{
			name: "dep1",
			ping: func(ctx context.Context) error {
				// Verify the context has tracing info
				span := oteltrace.SpanFromContext(ctx)
				require.True(t, span.SpanContext().IsValid())
				return nil
			},
		},
		&mockDep{
			name: "dep2",
			ping: func(ctx context.Context) error {
				// Verify the context has tracing info
				span := oteltrace.SpanFromContext(ctx)
				require.True(t, span.SpanContext().IsValid())
				return fmt.Errorf("dep2 error")
			},
		},
	}

	checker := NewChecker(deps...)
	_, success := checker.Check(parentCtx)
	require.False(t, success) // Should be false due to dep2 error

	// Force spans to be exported
	err := tp.ForceFlush(context.Background())
	require.NoError(t, err)

	// Verify spans were created
	spans := exporter.GetSpans()
	require.Len(t, spans, 2) // 2 dependency spans

	// Find dependency spans
	var dep1Span, dep2Span tracetest.SpanStub
	for _, span := range spans {
		switch span.Name {
		case "health.ping.dep1":
			dep1Span = span
		case "health.ping.dep2":
			dep2Span = span
		}
	}

	// Verify dependency spans exist and have correct parent
	require.NotEmpty(t, dep1Span.Name)
	require.NotEmpty(t, dep2Span.Name)
	require.Equal(t, parentSpan.SpanContext().TraceID(), dep1Span.SpanContext.TraceID())
	require.Equal(t, parentSpan.SpanContext().TraceID(), dep2Span.SpanContext.TraceID())
	require.Equal(t, parentSpan.SpanContext().SpanID(), dep1Span.Parent.SpanID())
	require.Equal(t, parentSpan.SpanContext().SpanID(), dep2Span.Parent.SpanID())

	// Verify span kinds
	require.Equal(t, oteltrace.SpanKindClient, dep1Span.SpanKind)
	require.Equal(t, oteltrace.SpanKindClient, dep2Span.SpanKind)

	// Verify span names & also the `name` attribute for UI friendliness
	var dep1NameAttr, dep2NameAttr string
	for _, attr := range dep1Span.Attributes {
		if attr.Key == "name" {
			dep1NameAttr = attr.Value.AsString()
			break
		}
	}
	for _, attr := range dep2Span.Attributes {
		if attr.Key == "name" {
			dep2NameAttr = attr.Value.AsString()
			break
		}
	}
	require.Equal(t, dep1Span.Name, dep1NameAttr)
	require.Equal(t, dep2Span.Name, dep2NameAttr)

	// Verify error handling on dep2 span
	require.True(t, len(dep2Span.Events) > 0) // Should have error event
	require.Contains(t, dep2Span.Status.Description, "ping failed")
}

type mockDep struct {
	name string
	ping func(ctx context.Context) error
}

func (m *mockDep) Name() string                   { return m.name }
func (m *mockDep) Ping(ctx context.Context) error { return m.ping(ctx) }

func singleHealthyDep(name string) []Pinger {
	return []Pinger{&mockDep{
		name: name,
		ping: func(ctx context.Context) error {
			return nil
		},
	}}
}

func singleUnhealthyDep(name string, err error) []Pinger {
	return []Pinger{&mockDep{
		name: name,
		ping: func(ctx context.Context) error {
			return err
		},
	}}
}

func multipleHealthyDeps(names ...string) []Pinger {
	deps := make([]Pinger, len(names))
	for i, name := range names {
		deps[i] = &mockDep{
			name: name,
			ping: func(ctx context.Context) error {
				return nil
			},
		}
	}
	return deps
}

func multipleUnhealthyDeps(err error, names ...string) []Pinger {
	deps := make([]Pinger, len(names))
	for i, name := range names {
		deps[i] = &mockDep{
			name: name,
			ping: func(ctx context.Context) error {
				return nil
			},
		}
	}
	deps[len(deps)-1].(*mockDep).ping = func(ctx context.Context) error {
		return err
	}
	return deps
}
