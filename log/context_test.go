package log

import (
	"context"
	"testing"
)

func TestDebugEnabled(t *testing.T) {
	ctx := Context(context.Background())
	if DebugEnabled(ctx) {
		t.Errorf("expected debug logs to be disabled")
	}
	ctx = Context(ctx, WithDebug())
	if !DebugEnabled(ctx) {
		t.Errorf("expected debug logs to be enabled")
	}
}
