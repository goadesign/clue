package debug

import (
	"context"
	"fmt"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := defaultOptions()
	if fmt.Sprintf("%p", opts.format) != fmt.Sprintf("%p", FormatJSON) {
		t.Errorf("got format %p, expected %p", opts.format, FormatJSON)
	}
	if opts.maxsize != DefaultMaxSize {
		t.Errorf("got maxsize %d, expected %d", opts.maxsize, DefaultMaxSize)
	}
}

func TestWithFormat(t *testing.T) {
	opts := defaultOptions()
	WithFormat(FormatJSON)(opts)
	if fmt.Sprintf("%p", opts.format) != fmt.Sprintf("%p", FormatJSON) {
		t.Errorf("got format %p, expected %p", opts.format, FormatJSON)
	}
}

func TestWithMaxSize(t *testing.T) {
	opts := defaultOptions()
	WithMaxSize(10)(opts)
	if opts.maxsize != 10 {
		t.Errorf("got maxsize %d, expected 10", opts.maxsize)
	}
}

func TestFormatJSON(t *testing.T) {
	ctx := context.Background()
	js := FormatJSON(ctx, map[string]string{"foo": "bar"})
	if js != `{"foo":"bar"}` {
		t.Errorf("got %q, expected %q", js, `{"foo":"bar"}`)
	}
	js = FormatJSON(ctx, make(chan int))
	if js != `<invalid: json: unsupported type: chan int>` {
		t.Errorf("got %q, expected %q", js, `<invalid: json: unsupported type: chan int>`)
	}

}
