package log

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestDefaultOptions(t *testing.T) {
	opts := defaultOptions()
	if fmt.Sprintf("%p", opts.disableBuffering) != fmt.Sprintf("%p", IsTracing) {
		t.Errorf("got disable buffering %p, expected %p", &opts.disableBuffering, IsTracing)
	}
	if opts.debug {
		t.Errorf("expected debug to be disabled")
	}
	if opts.w != os.Stdout {
		t.Errorf("got output %p, expected os.Stdout", opts.w)
	}
	if fmt.Sprintf("%p", opts.format) != fmt.Sprintf("%p", FormatText) {
		t.Errorf("got format %p, expected %p", opts.format, FormatText)
	}
	if opts.maxsize != DefaultMaxSize {
		t.Errorf("got maxsize %d, expected %d", opts.maxsize, DefaultMaxSize)
	}
}

func TestWithDisableBuffering(t *testing.T) {
	opts := defaultOptions()
	disable := func(ctx context.Context) bool { return true }
	WithDisableBuffering(disable)(opts)
	if fmt.Sprintf("%p", opts.disableBuffering) != fmt.Sprintf("%p", disable) {
		t.Errorf("got disable buffering %p, expected %p", opts.disableBuffering, disable)
	}
}

func TestWithDebug(t *testing.T) {
	opts := defaultOptions()
	WithDebug()(opts)
	if !opts.debug {
		t.Errorf("expected debug to be enabled")
	}
}

func TestWithOutput(t *testing.T) {
	opts := defaultOptions()
	w := ioutil.Discard
	WithOutput(w)(opts)
	if fmt.Sprintf("%p", opts.w) != fmt.Sprintf("%p", w) {
		t.Errorf("got output %p, expected %p", opts.w, w)
	}
}

func TestWithFormat(t *testing.T) {
	opts := defaultOptions()
	WithFormat(FormatJSON)(opts)
	if fmt.Sprintf("%p", opts.format) != fmt.Sprintf("%p", FormatJSON) {
		t.Errorf("got format %p, expected %p", opts.format, FormatJSON)
	}
}

func TestWithKeyValue(t *testing.T) {
	opts := defaultOptions()
	WithKeyValue("key", "value")(opts)
	if len(opts.keyvals) != 2 {
		t.Fatalf("got %d keyvals, expected 2", len(opts.keyvals))
	}
	if opts.keyvals[0] != "key" {
		t.Errorf("got key %q, expected \"key\"", opts.keyvals[0])
	}
	if opts.keyvals[1] != "value" {
		t.Errorf("got value %q, expected \"value\"", opts.keyvals[1])
	}
}

func TestWithMaxSize(t *testing.T) {
	opts := defaultOptions()
	WithMaxSize(10)(opts)
	if opts.maxsize != 10 {
		t.Errorf("got maxsize %d, expected 10", opts.maxsize)
	}
}

func TestIsTracing(t *testing.T) {
	if IsTracing(context.Background()) {
		t.Errorf("expected IsTracing to return false")
	}

	exp, _ := stdouttrace.New(stdouttrace.WithWriter(ioutil.Discard))
	tp := trace.NewTracerProvider(trace.WithBatcher(exp))
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)
	ctx, span := otel.Tracer("test").Start(context.Background(), "test")
	defer span.End()
	if !IsTracing(ctx) {
		t.Errorf("expected IsTracing to return true")
	}
}

func TestIsTerminal(t *testing.T) {
	if IsTerminal() {
		t.Errorf("expected IsTerminal to return false")
	}
}
