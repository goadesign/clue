package trace

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"goa.design/clue/internal/testsvc"
)

func TestHTTP(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	ctx := testContext(provider)
	cli, stop := testsvc.SetupHTTP(t,
		testsvc.WithHTTPMiddleware(HTTP(ctx)),
		testsvc.WithHTTPFunc(addEventUnaryMethod))
	_, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{})
	assert.NoError(t, err)
	stop()
	spans := exporter.GetSpans()
	require.Len(t, spans, 1)
	assert.Equal(t, "test", spans[0].Name)
}

func TestClient(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	ctx := testContext(provider)
	c := http.Client{Transport: Client(ctx, http.DefaultTransport)}
	otelt, ok := c.Transport.(*otelhttp.Transport)
	assert.True(t, ok, "got %T, want %T", c.Transport, otelt)
}
