package log

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"goa.design/goa/v3/middleware"
)

func TestInit(t *testing.T) {
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()

	endpoint := func(ctx context.Context, req interface{}) (interface{}, error) {
		Print(ctx, "hello world", "key1", "value1", "key2", "value2")
		return nil, nil
	}
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatJSON))
	requestIDCtx := context.WithValue(context.Background(), middleware.RequestIDKey, "request-id")

	Init(ctx)(endpoint)(requestIDCtx, nil)

	expected := fmt.Sprintf("{%s,%s,%s,%s,%s,%s}\n",
		`"level":"INFO"`,
		`"time":"2022-01-09T20:29:45Z"`,
		`"msg":"hello world"`,
		`"request_id":"request-id"`,
		`"key1":"value1"`,
		`"key2":"value2"`)

	if buf.String() != expected {
		t.Errorf("got %s, want %s", buf.String(), expected)
	}
}
