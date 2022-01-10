package log

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestSetContext(t *testing.T) {
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()

	endpoint := func(ctx context.Context, req interface{}) (interface{}, error) {
		Print(ctx, "hello world", "key1", "value1", "key2", "value2")
		return nil, nil
	}
	var buf bytes.Buffer
	endpoint = SetContext(WithOutput(&buf), WithFormat(FormatJSON))(endpoint)

	endpoint(context.Background(), nil)

	if buf.String() != `{"level":"INFO","time":"2022-01-09T20:29:45Z","msg":"hello world","key1":"value1","key2":"value2"}`+"\n" {
		t.Errorf("got %s, want %s", buf.String(), `{"level":"INFO","time":"2022-01-09T20:29:45Z","msg":"hello world","key1":"value1","key2":"value2"}`)
	}
}
