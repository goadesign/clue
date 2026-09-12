package log

import (
	"bytes"
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

func TestWithIsolatesOptions(t *testing.T) {
	var original, changed bytes.Buffer
	base := Context(context.Background(), WithDebug(), WithOutputs(Output{Writer: &original, Format: testFormat}))
	copy := Context(With(base), WithNoDebug(), WithOutput(&changed), WithMaxSize(4))
	Debugf(base, "original")
	Debugf(copy, "hidden")
	Printf(copy, "changed")
	assert.Equal(t, "original", original.String())
	assert.Equal(t, "chan ... <clue/log.truncated>", changed.String())
	assert.True(t, DebugEnabled(base))
	assert.False(t, DebugEnabled(copy))
}
