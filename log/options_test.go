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
	if assert.Len(t, opts.outputs, 1) {
		assert.Equal(t, os.Stdout, opts.outputs[0].Writer)
		assert.Equal(t, fmt.Sprintf("%p", opts.outputs[0].Format), fmt.Sprintf("%p", FormatText))
	}
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

func TestWithOutputs(t *testing.T) {
	opts := defaultOptions()
	w := io.Discard
	WithOutputs(Output{Writer: w, Format: FormatJSON})(opts)
	if assert.Len(t, opts.outputs, 1) {
		assert.Equal(t, w, opts.outputs[0].Writer)
		assert.Equal(t, fmt.Sprintf("%p", opts.outputs[0].Format), fmt.Sprintf("%p", FormatJSON))
	}
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
