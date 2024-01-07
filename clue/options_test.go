package clue

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"goa.design/clue/log"
)

func TestOptions(t *testing.T) {
	ctx := log.Context(context.Background())
	cases := []struct {
		name    string
		options []Option
		want    *options
	}{
		{
			name: "default",
			want: defaultOptions(ctx),
		},
		{
			name:    "with reader interval",
			options: []Option{WithReaderInterval(10)},
			want: func() *options {
				o := defaultOptions(ctx)
				o.readerInterval = 10
				return o
			}(),
		},
		{
			name:    "with max sampling rate",
			options: []Option{WithMaxSamplingRate(10)},
			want: func() *options {
				o := defaultOptions(ctx)
				o.maxSamplingRate = 10
				return o
			}(),
		},
		{
			name:    "with sample size",
			options: []Option{WithSampleSize(10)},
			want: func() *options {
				o := defaultOptions(ctx)
				o.sampleSize = 10
				return o
			}(),
		},
		{
			name:    "with propagator",
			options: []Option{WithPropagators(nil)},
			want: func() *options {
				o := defaultOptions(ctx)
				o.propagators = nil
				return o
			}(),
		},
		{
			name:    "with error handler",
			options: []Option{WithErrorHandler(dummyErrorHandler{})},
			want: func() *options {
				o := defaultOptions(ctx)
				o.errorHandler = dummyErrorHandler{}
				return o
			}(),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := defaultOptions(ctx)
			for _, o := range c.options {
				o(got)
			}
			assert.Equal(t, c.want, got)
		})
	}
}
