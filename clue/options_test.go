package clue

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"

	"goa.design/clue/log"
)

func TestOptions(t *testing.T) {
	ctx := log.Context(context.Background())
	cases := []struct {
		name   string
		option Option
		want   func(*options) // mutate default options
	}{
		{
			name: "default",
		},
		{
			name:   "with reader interval",
			option: WithReaderInterval(1000),
			want:   func(o *options) { o.readerInterval = 1000 },
		},
		{
			name:   "with max sampling rate",
			option: WithMaxSamplingRate(1000),
			want:   func(o *options) { o.maxSamplingRate = 1000 },
		},
		{
			name:   "with sample size",
			option: WithSampleSize(1000),
			want:   func(o *options) { o.sampleSize = 1000 },
		},
		{
			name:   "with propagator",
			option: WithPropagators(propagation.TraceContext{}),
			want:   func(o *options) { o.propagators = propagation.TraceContext{} },
		},
		{
			name:   "with resource",
			option: WithResource(resource.NewSchemaless(attribute.String("key", "value"))),
			want:   func(o *options) { o.resource = resource.NewSchemaless(attribute.String("key", "value")) },
		},
		{
			name:   "with error handler",
			option: WithErrorHandler(dummyErrorHandler{}),
			want:   func(o *options) { o.errorHandler = dummyErrorHandler{} },
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := defaultOptions(ctx)
			if c.option != nil {
				c.option(got)
			}
			want := defaultOptions(ctx)
			if c.want != nil {
				c.want(want)
			}
			assert.Equal(t, want.maxSamplingRate, got.maxSamplingRate)
			assert.Equal(t, want.sampleSize, got.sampleSize)
			assert.Equal(t, want.readerInterval, got.readerInterval)
			assert.Equal(t, want.propagators, got.propagators)
			assert.Equal(t, want.resource, got.resource)
			assert.IsType(t, want.errorHandler, got.errorHandler)
		})
	}
}
