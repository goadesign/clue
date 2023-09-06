package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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
	assert.NoError(t, WithMaxSamplingRate(3)(ctx, options))
	if options.maxSamplingRate != 3 {
		t.Errorf("got %d sampling rate, want 3", options.maxSamplingRate)
	}
	assert.NoError(t, WithSampleSize(20)(ctx, options))
	if options.sampleSize != 20 {
		t.Errorf("got %d sample size, want 20", options.sampleSize)
	}
	assert.NoError(t, WithDisabled()(ctx, options))
	if !options.disabled {
		t.Error("expected disabled to be true")
	}
	assert.NoError(t, WithExporter(tracetest.NewInMemoryExporter())(ctx, options))
	if options.exporter == nil {
		t.Error("got nil exporter, want non-nil")
	}
	assert.NoError(t, WithResource(&resource.Resource{})(ctx, options))
	if options.resource == nil {
		t.Error("got nil resource, want non-nil")
	}
	assert.NoError(t, WithParentSamplerOptions(sdktrace.WithRemoteParentSampled(nil))(ctx, options))
	if total := len(options.parentSamplerOptions); total != 1 {
		t.Errorf("got %d parent sampler options, expected 1", total)
	}
}
