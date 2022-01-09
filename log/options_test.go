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
