package tracing

import (
	"context"
	"testing"

	"github.com/crossnokaye/micro/internal/testsvc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"goa.design/goa/v3/http/middleware"
)

func TestMiddleware(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	cli, stop := testsvc.SetupHTTP(t,
		testsvc.WithHTTPMiddleware(middleware.RequestID(), Middleware("test", provider)),
		testsvc.WithHTTPFunc(noopUnaryMethod))
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
}

func TestMiddlewareNoRequestID(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	cli, stop := testsvc.SetupHTTP(t,
		testsvc.WithHTTPMiddleware(Middleware("test", provider)),
		testsvc.WithHTTPFunc(noopUnaryMethod))
	if _, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	stop()
	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
}
