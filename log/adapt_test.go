package log

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestAdapt(t *testing.T) {
	restore := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = restore }()

	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf))
	logger := Adapt(ctx)
	logger.Log("msg", "hello")
	want := "INFO[2022-01-09T20:29:45Z] msg=hello\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
}
