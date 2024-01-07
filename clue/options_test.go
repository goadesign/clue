package clue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyErrorHandler struct{}

func (dummyErrorHandler) Handle(error) {}

func TestOptions(t *testing.T) {
	cases := []struct {
		name    string
		options []Option
		want    *options
	}{
		{
			name: "default",
			want: defaultOptions(),
		},
		{
			name:    "with reader interval",
			options: []Option{WithReaderInterval(10)},
			want: func() *options {
				o := defaultOptions()
				o.readerInterval = 10
				return o
			}(),
		},
		{
			name:    "with max sampling rate",
			options: []Option{WithMaxSamplingRate(10)},
			want: func() *options {
				o := defaultOptions()
				o.maxSamplingRate = 10
				return o
			}(),
		},
		{
			name:    "with sample size",
			options: []Option{WithSampleSize(10)},
			want: func() *options {
				o := defaultOptions()
				o.sampleSize = 10
				return o
			}(),
		},
		{
			name:    "with propagator",
			options: []Option{WithPropagators(nil)},
			want: func() *options {
				o := defaultOptions()
				o.propagators = nil
				return o
			}(),
		},
		{
			name:    "with error handler",
			options: []Option{WithErrorHandler(dummyErrorHandler{})},
			want: func() *options {
				o := defaultOptions()
				o.errorHandler = dummyErrorHandler{}
				return o
			}(),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := defaultOptions()
			for _, o := range c.options {
				o(got)
			}
			assert.Equal(t, c.want, got)
		})
	}
}
