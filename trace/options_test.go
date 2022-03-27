package trace

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestOptions(t *testing.T) {
	ctx := context.Background()

	options := defaultOptions()
	if options.maxSamplingRate != 2 {
		t.Errorf("got %d, want 2", options.maxSamplingRate)
	}
	if options.sampleSize != 10 {
		t.Errorf("got %d, want 10", options.sampleSize)
	}
	WithMaxSamplingRate(3)(ctx, options)
	if options.maxSamplingRate != 3 {
		t.Errorf("got %d sampling rate, want 3", options.maxSamplingRate)
	}
	WithSampleSize(20)(ctx, options)
	if options.sampleSize != 20 {
		t.Errorf("got %d sample size, want 20", options.sampleSize)
	}
	WithDisabled()(ctx, options)
	if !options.disabled {
		t.Error("expected disabled to be true")
	}
	WithExporter(tracetest.NewInMemoryExporter())(ctx, options)
	if options.exporter == nil {
		t.Error("got nil exporter, want non-nil")
	}
}
