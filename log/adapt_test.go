package log

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/aws/smithy-go/logging"
	"github.com/stretchr/testify/assert"
)

func TestAsGoaMiddlwareLogger(t *testing.T) {
	restore := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = restore }()

	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf))
	logger := AsGoaMiddlewareLogger(ctx)
	assert.NoError(t, logger.Log("msg", "hello world"))
	want := "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, buf.String(), want)
}

func TestAsStdLogger(t *testing.T) {
	restore := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = restore }()

	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf))
	logger := AsStdLogger(ctx)

	logger.Print("hello world")
	want := "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, buf.String(), want)

	buf.Reset()
	logger.Println("hello world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\\n\"\n"
	assert.Equal(t, buf.String(), want)

	buf.Reset()
	logger.Printf("hello %s", "world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, buf.String(), want)

	func() {
		buf.Reset()
		var msg string
		defer func() { msg = recover().(string) }()
		logger.Panic("hello world")
		want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
		assert.Equal(t, buf.String(), want)
		assert.Equal(t, msg, "hello world")
	}()

	func() {
		buf.Reset()
		var msg string
		defer func() { msg = recover().(string) }()
		logger.Panicf("hello %s", "world")
		want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
		assert.Equal(t, buf.String(), want)
		assert.Equal(t, msg, "hello world")
	}()

	func() {
		buf.Reset()
		var msg string
		defer func() { msg = recover().(string) }()
		logger.Panicln("hello world")
		want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\\n\"\n"
		assert.Equal(t, buf.String(), want)
		assert.Equal(t, msg, "hello world")
	}()

	osExitFunc := osExit
	var exited int
	osExit = func(code int) {
		exited = code
	}
	defer func() { osExit = osExitFunc }()

	buf.Reset()
	logger.Fatal("hello world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, buf.String(), want)
	assert.Equal(t, exited, 1)

	exited = 0
	buf.Reset()
	logger.Fatalf("hello %s", "world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, buf.String(), want)
	assert.Equal(t, exited, 1)

	exited = 0
	buf.Reset()
	logger.Fatalln("hello world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\\n\"\n"
	assert.Equal(t, buf.String(), want)
	assert.Equal(t, exited, 1)
}

type ctxkey string

func TestAsAWSLogger(t *testing.T) {
	restore := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = restore }()
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithDebug())
	var logger logging.Logger = AsAWSLogger(ctx)

	logger.Logf(logging.Classification("INFO"), "hello %s", "world")
	want := "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, buf.String(), want)

	buf.Reset()
	logger.Logf(logging.Classification("DEBUG"), "hello world")
	want = "time=2022-01-09T20:29:45Z level=debug msg=\"hello world\"\n"
	assert.Equal(t, buf.String(), want)

	buf.Reset()
	key := ctxkey("key")
	pctx := context.WithValue(context.Background(), key, "small")
	logger = logger.(logging.ContextLogger).WithContext(pctx)
	logger.Logf(logging.Classification("INFO"), "hello %v world", logger.(*AWSLogger).Context.Value(key))
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello small world\"\n"
	assert.Equal(t, buf.String(), want)
}
