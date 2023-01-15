package debug

import (
	"context"
	"fmt"
	"testing"
)

func TestDefaultLogPayloadsOptions(t *testing.T) {
	opts := defaultLogPayloadsOptions()
	if fmt.Sprintf("%p", opts.format) != fmt.Sprintf("%p", FormatJSON) {
		t.Errorf("got format %p, expected %p", opts.format, FormatJSON)
	}
	if opts.maxsize != DefaultMaxSize {
		t.Errorf("got maxsize %d, expected %d", opts.maxsize, DefaultMaxSize)
	}
}

func TestDefaultDebugLogEnablerOptions(t *testing.T) {
	opts := defaultDebugLogEnablerOptions()
	if opts.path != "debug" {
		t.Errorf("got prefix %q, expected %q", opts.path, "debug")
	}
	if opts.query != "debug-logs" {
		t.Errorf("got query %q, expected %q", opts.query, "debug-logs")
	}
	if opts.onval != "on" {
		t.Errorf("got onval %q, expected %q", opts.onval, "on")
	}
	if opts.offval != "off" {
		t.Errorf("got offval %q, expected %q", opts.offval, "off")
	}
}

func TestDefaultPprofOptions(t *testing.T) {
	opts := defaultPprofOptions()
	if opts.prefix != "/debug/pprof/" {
		t.Errorf("got prefix %q, expected %q", opts.prefix, "/debug/pprof/")
	}
}

func TestWithFormat(t *testing.T) {
	opts := defaultLogPayloadsOptions()
	WithFormat(FormatJSON)(opts)
	if fmt.Sprintf("%p", opts.format) != fmt.Sprintf("%p", FormatJSON) {
		t.Errorf("got format %p, expected %p", opts.format, FormatJSON)
	}
}

func TestWithMaxSize(t *testing.T) {
	opts := defaultLogPayloadsOptions()
	WithMaxSize(10)(opts)
	if opts.maxsize != 10 {
		t.Errorf("got maxsize %d, expected 10", opts.maxsize)
	}
}

func TestWithClient(t *testing.T) {
	opts := defaultLogPayloadsOptions()
	WithClient()(opts)
	if !opts.client {
		t.Errorf("got client %v, expected true", opts.client)
	}
}

func TestWithPrefix(t *testing.T) {
	opts := defaultDebugLogEnablerOptions()
	WithPath("foo")(opts)
	if opts.path != "foo" {
		t.Errorf("got prefix %q, expected %q", opts.path, "foo")
	}
}

func TestWithQuery(t *testing.T) {
	opts := defaultDebugLogEnablerOptions()
	WithQuery("foo")(opts)
	if opts.query != "foo" {
		t.Errorf("got query %q, expected %q", opts.query, "foo")
	}
}

func TestWithOnValue(t *testing.T) {
	opts := defaultDebugLogEnablerOptions()
	WithOnValue("foo")(opts)
	if opts.onval != "foo" {
		t.Errorf("got onval %q, expected %q", opts.onval, "foo")
	}
}

func TestWithOffValue(t *testing.T) {
	opts := defaultDebugLogEnablerOptions()
	WithOffValue("foo")(opts)
	if opts.offval != "foo" {
		t.Errorf("got offval %q, expected %q", opts.offval, "foo")
	}
}

func TestWithPprofPrefix(t *testing.T) {
	opts := defaultPprofOptions()
	WithPrefix("foo")(opts)
	if opts.prefix != "foo" {
		t.Errorf("got pprofPrefix %q, expected %q", opts.prefix, "foo")
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
