package log

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"goa.design/goa/v3/middleware"
)

func TestHTTP(t *testing.T) {
	now := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = now }()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Print(req.Context(), "hello world", "key1", "value1", "key2", "value2")
	})
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatJSON))

	handler = HTTP(ctx)(handler)

	requestIDCtx := context.WithValue(ctx, middleware.RequestIDKey, "request-id")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req = req.WithContext(requestIDCtx)

	handler.ServeHTTP(nil, req)

	expected := fmt.Sprintf("{%s,%s,%s,%s,%s,%s}\n",
		`"time":"2022-01-09T20:29:45Z"`,
		`"level":"info"`,
		`"msg":"hello world"`,
		`"requestID":"request-id"`,
		`"key1":"value1"`,
		`"key2":"value2"`)

	if buf.String() != expected {
		t.Errorf("got %s, want %s", buf.String(), expected)
	}
}
