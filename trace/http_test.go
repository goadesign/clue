package trace

import (
	"context"
	"net/http"
	"testing"

	"github.com/crossnokaye/micro/internal/testsvc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"goa.design/goa/v3/http/middleware"
)

func TestHTTP(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	ctx := withProvider(context.Background(), provider)
	cli, stop := testsvc.SetupHTTP(t,
		testsvc.WithHTTPMiddleware(middleware.RequestID(), HTTP(ctx, "test")),
		testsvc.WithHTTPFunc(addEventUnaryMethod))
	if _, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	stop()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
	found := false
	for _, att := range spans[0].Attributes {
		if att.Key == AttributeRequestID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("request ID not in span attributes")
	}
	events := spans[0].Events
	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}
	if events[0].Name != "unary method" {
		t.Errorf("unexpected event name: %s", events[0].Name)
	}
}

func TestHTTPRequestID(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	ctx := withProvider(context.Background(), provider)
	cli, stop := testsvc.SetupHTTP(t,
		testsvc.WithHTTPMiddleware(HTTP(ctx, "test")),
		testsvc.WithHTTPFunc(addEventUnaryMethod))
	if _, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	stop()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
}

func TestTrace(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	ctx := withProvider(context.Background(), provider)
	c := http.Client{Transport: Client(ctx, http.DefaultTransport)}
	otelt, ok := c.Transport.(*otelhttp.Transport)
	if !ok {
		t.Errorf("got %T, want %T", c.Transport, otelt)
	}
}
