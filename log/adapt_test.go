package log

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/smithy-go/logging"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
)

func TestAsGoaMiddlwareLogger(t *testing.T) {
	restore := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = restore }()

	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutputs(Output{Writer: &buf, Format: FormatText}))
	logger := AsGoaMiddlewareLogger(ctx)
	assert.NoError(t, logger.Log("msg", "hello world"))
	want := "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, want, buf.String())
}

func TestAsStdLogger(t *testing.T) {
	restore := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = restore }()

	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutputs(Output{Writer: &buf, Format: FormatText}))
	logger := AsStdLogger(ctx)

	logger.Print("hello world")
	want := "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, want, buf.String())

	buf.Reset()
	logger.Println("hello world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\n\"\n"
	assert.Equal(t, want, buf.String())

	buf.Reset()
	logger.Printf("hello %s", "world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, want, buf.String())

	func() {
		buf.Reset()
		var msg string
		defer func() { msg = recover().(string) }()
		logger.Panic("hello world")
		want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
		assert.Equal(t, want, buf.String())
		assert.Equal(t, msg, "hello world")
	}()

	func() {
		buf.Reset()
		var msg string
		defer func() { msg = recover().(string) }()
		logger.Panicf("hello %s", "world")
		want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
		assert.Equal(t, want, buf.String())
		assert.Equal(t, msg, "hello world")
	}()

	func() {
		buf.Reset()
		var msg string
		defer func() { msg = recover().(string) }()
		logger.Panicln("hello world")
		want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\\n\"\n"
		assert.Equal(t, want, buf.String())
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
	assert.Equal(t, want, buf.String())
	assert.Equal(t, exited, 1)

	exited = 0
	buf.Reset()
	logger.Fatalf("hello %s", "world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, want, buf.String())
	assert.Equal(t, exited, 1)

	exited = 0
	buf.Reset()
	logger.Fatalln("hello world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\n\"\n"
	assert.Equal(t, want, buf.String())
	assert.Equal(t, exited, 1)
}

type ctxkey string

func TestAsAWSLogger(t *testing.T) {
	restore := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = restore }()
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutputs(Output{Writer: &buf, Format: FormatText}), WithDebug())
	var logger logging.Logger = AsAWSLogger(ctx)

	logger.Logf(logging.Classification("INFO"), "hello %s", "world")
	want := "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	assert.Equal(t, want, buf.String())

	buf.Reset()
	logger.Logf(logging.Classification("DEBUG"), "hello world")
	want = "time=2022-01-09T20:29:45Z level=debug msg=\"hello world\"\n"
	assert.Equal(t, want, buf.String())

	buf.Reset()
	key := ctxkey("key")
	pctx := context.WithValue(context.Background(), key, "small")
	logger = logger.(logging.ContextLogger).WithContext(pctx)
	logger.Logf(logging.Classification("INFO"), "hello %v world", pctx.Value(key))
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello small world\"\n"
	assert.Equal(t, want, buf.String())
}

func TestToLogrSink(t *testing.T) {
	restore := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = restore }()
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutputs(Output{Writer: &buf, Format: FormatText}), WithDebug())
	var sink logr.LogSink = ToLogrSink(ctx)

	sink.Init(logr.RuntimeInfo{})
	assert.True(t, sink.Enabled(0))
	assert.True(t, sink.Enabled(1))
	assert.True(t, sink.Enabled(2))
	assert.True(t, sink.Enabled(3))

	msg := "hello world"
	expected := "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	expecteddebug := "time=2022-01-09T20:29:45Z level=debug msg=\"hello world\"\n"
	expectederr := "time=2022-01-09T20:29:45Z level=error err=error msg=\"hello world\"\n"
	empty := ""

	logger := logr.New(sink)
	logger.Info(msg)
	assert.Equal(t, expected, buf.String())

	buf.Reset()
	logger.V(1).Info(msg)
	assert.Equal(t, expecteddebug, buf.String())

	buf.Reset()
	logger.Error(errors.New("error"), msg)
	assert.Equal(t, expectederr, buf.String())

	ctx = Context(context.Background(), WithOutputs(Output{Writer: &buf, Format: FormatText}))
	sink = ToLogrSink(ctx)
	logger = logr.New(sink)

	buf.Reset()
	logger.Info(msg)
	assert.Equal(t, empty, buf.String())

	FlushAndDisableBuffering(ctx)
	buf.Reset()
	logger.Info(msg)
	assert.Equal(t, expected, buf.String())

	buf.Reset()
	logger.V(1).Info(msg)
	assert.Equal(t, empty, buf.String())

	sink = sink.WithValues("key", "value")
	expectedWithValues := "time=2022-01-09T20:29:45Z level=info key=value msg=\"hello world\"\n"
	logger = logr.New(sink)
	buf.Reset()
	logger.Info(msg)
	assert.Equal(t, expectedWithValues, buf.String())

	sink = sink.WithName("name")
	expectedWithName := "time=2022-01-09T20:29:45Z level=info key=value log=name msg=\"hello world\"\n"
	logger = logr.New(sink)
	buf.Reset()
	logger.Info(msg)
	assert.Equal(t, expectedWithName, buf.String())
}
