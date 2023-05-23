package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugEnabled(t *testing.T) {
	ctx := Context(context.Background())
	assert.False(t, DebugEnabled(ctx), "expected debug logs to be disabled")
	ctx = Context(ctx, WithDebug())
	assert.True(t, DebugEnabled(ctx), "expected debug logs to be enabled")
}
