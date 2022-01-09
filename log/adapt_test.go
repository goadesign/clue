package log

import (
	"bytes"
	"context"
	"testing"
)

func TestAdapt(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf))
	logger := Adapt(ctx)
	logger.Log("msg", "hello")
	if buf.String() != "[INFO] [msg=hello]\n" {
		t.Errorf("got %q, want %q", buf.String(), "hello")
	}
}
