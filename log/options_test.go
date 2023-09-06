package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestDefaultOptions(t *testing.T) {
	opts := defaultOptions()
	assert.Equal(t, fmt.Sprintf("%p", opts.disableBuffering), fmt.Sprintf("%p", IsTracing))
	assert.False(t, opts.debug)
	assert.Equal(t, opts.w, os.Stdout)
	assert.Equal(t, fmt.Sprintf("%p", opts.format), fmt.Sprintf("%p", FormatText))
	assert.Equal(t, opts.maxsize, DefaultMaxSize)
}

func TestWithDisableBuffering(t *testing.T) {
	opts := defaultOptions()
	disable := func(ctx context.Context) bool { return true }
	WithDisableBuffering(disable)(opts)
	assert.Equal(t, fmt.Sprintf("%p", opts.disableBuffering), fmt.Sprintf("%p", disable))
}

func TestWithDebug(t *testing.T) {
	opts := defaultOptions()
	WithDebug()(opts)
	assert.True(t, opts.debug)
}

func TestWithNoDebug(t *testing.T) {
	opts := defaultOptions()
	WithDebug()(opts)
	WithNoDebug()(opts)
	assert.False(t, opts.debug)
}

func TestWithOutput(t *testing.T) {
	opts := defaultOptions()
	w := io.Discard
	WithOutput(w)(opts)
	assert.Equal(t, opts.w, w)
}

func TestWithFormat(t *testing.T) {
	opts := defaultOptions()
	WithFormat(FormatJSON)(opts)
	assert.Equal(t, fmt.Sprintf("%p", opts.format), fmt.Sprintf("%p", FormatJSON))
}

func TestWithMaxSize(t *testing.T) {
	opts := defaultOptions()
	WithMaxSize(10)(opts)
	assert.Equal(t, opts.maxsize, 10)
}

func TestIsTracing(t *testing.T) {
	if IsTracing(context.Background()) {
		t.Errorf("expected IsTracing to return false")
	}
	exp, _ := stdouttrace.New(stdouttrace.WithWriter(io.Discard))
	tp := trace.NewTracerProvider(trace.WithBatcher(exp))
	defer tp.Shutdown(context.Background()) // nolint:errcheck
	otel.SetTracerProvider(tp)
	ctx, span := otel.Tracer("test").Start(context.Background(), "test")
	defer span.End()
	assert.True(t, IsTracing(ctx))
}

func TestIsTerminal(t *testing.T) {
	if IsTerminal() {
		t.Errorf("expected IsTerminal to return false")
	}
}
