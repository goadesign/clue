package log

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Buffer is a goroutine safe bytes.Buffer
type Buffer struct {
	buffer bytes.Buffer
	mutex  sync.Mutex
}

const (
	buffered = "buffered"
	printed  = "printed"
	ignored  = "ignored"
)

const ServiceName = "service"

func init() {
	// Mock time.Now for examples and tests
	timeNow = func() time.Time {
		return time.Date(2022, time.February, 22, 17, 0, 0, 0, time.UTC)
	}
}

func (s *Buffer) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.Write(p)
}

func (s *Buffer) String() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.String()
}

func ExamplePrintf() {
	ctx := Context(context.Background())
	Printf(ctx, "hello %s", "world")
	// Output: time=2022-02-22T17:00:00Z level=info msg="hello world"
}

func ExamplePrint() {
	ctx := Context(context.Background())
	Print(ctx, KV{"hello", "world"})
	// Output: time=2022-02-22T17:00:00Z level=info hello=world
}

func ExampleErrorf() {
	ctx := Context(context.Background())
	err := errors.New("error")
	Info(ctx, KV{"hello", "world"})
	// No output at this point because Info log events are buffered.
	// The call to Errorf causes the buffered events to be flushed.
	fmt.Println("---")
	Errorf(ctx, err, "failure")
	// Output: ---
	// time=2022-02-22T17:00:00Z level=info hello=world
	// time=2022-02-22T17:00:00Z level=error err=error msg=failure
}

// Note: if the line number for the call to Infof below changes update the test
// accordingly.
func TestFileLocation(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(testFormat), WithFileLocation())
	Infof(ctx, buffered)
	require.Len(t, entries(ctx), 1)
	e := (entries(ctx))[0]
	require.Len(t, e.KeyVals, 2)
	assert.Equal(t, KV{"msg", buffered}, e.KeyVals[0])
	assert.Equal(t, KV{"file", "log/log_test.go:82"}, e.KeyVals[1])
}

func TestSeverity(t *testing.T) {
	var buf bytes.Buffer
	printSev := func(e *Entry) []byte {
		return []byte(e.Severity.String() + ":" + e.Severity.Code() + ":" + e.Severity.Color() + " ")
	}
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(printSev), WithDebug())
	Debugf(ctx, "")
	Infof(ctx, "")
	Errorf(ctx, nil, "")
	want := "debug:DEBG:\033[37m info:INFO:\033[34m error:ERRO:\033[1;31m "
	assert.Equal(t, want, buf.String())
	assert.Equal(t, "<INVALID>", Severity(0).String())
	assert.Equal(t, "<INVALID>", Severity(0).Code())
	assert.Empty(t, Severity(0).Color())
}

func TestBuffering(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(testFormat))

	// Buffering is enabled by default.
	Infof(ctx, buffered)
	require.Len(t, entries(ctx), 1)
	e := entries(ctx)[0]
	require.Len(t, e.KeyVals, 1)
	assert.Equal(t, KV{"msg", buffered}, e.KeyVals[0])

	// Printf does not buffer.
	Printf(ctx, printed)
	assert.Equal(t, printed, buf.String())

	// Flush flushes the buffer.
	FlushAndDisableBuffering(ctx)
	assert.Empty(t, entries(ctx))
	assert.Equal(t, printed+buffered, buf.String())

	// Buffering is disabled after flush.
	Infof(ctx, printed)
	assert.Empty(t, entries(ctx))
	assert.Equal(t, printed+buffered+printed, buf.String())

	// Flush is idempotent.
	FlushAndDisableBuffering(ctx)
	Infof(ctx, printed)
	assert.Empty(t, entries(ctx))
	assert.Equal(t, printed+buffered+printed+printed, buf.String())
}

func TestBufferingWithError(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(testFormat))
	err := fmt.Errorf("error")

	// Error flushes the buffer.
	Infof(ctx, buffered)
	Errorf(ctx, err, printed)
	assert.Empty(t, entries(ctx))
	expected := buffered + err.Error() + printed
	assert.Equal(t, expected, buf.String())

	// Buffering is disabled after error.
	Infof(ctx, printed)
	assert.Empty(t, entries(ctx))
	assert.Equal(t, expected+printed, buf.String())
}

type ctxTestKey int

const disableBufferingKey ctxTestKey = iota + 1

func TestBufferingWithDisableBufferingFunc(t *testing.T) {
	disableBuffering := func(ctx context.Context) bool {
		return ctx.Value(disableBufferingKey) == true
	}

	cases := []struct {
		name    string
		ctxFunc func(context.Context) context.Context
	}{
		{"with", func(ctx context.Context) context.Context { return With(ctx) }},
		{"context", func(ctx context.Context) context.Context { return Context(ctx) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := Context(context.Background(), WithOutput(&buf),
				WithFormat(testFormat), WithDisableBuffering(disableBuffering))

			Infof(ctx, buffered)
			assert.Len(t, entries(ctx), 1)

			ctx = tc.ctxFunc(context.WithValue(ctx, disableBufferingKey, true))
			Infof(ctx, printed)

			expected := buffered + printed
			assert.Equal(t, expected, buf.String())
		})
	}
}

func TestFatal(t *testing.T) {
	var exitCalled bool
	osExit = func(_ int) { exitCalled = true }
	defer func() { osExit = os.Exit }()
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(testFormat))
	err := fmt.Errorf("error")

	// Fatal flushes the buffer.
	Infof(ctx, buffered)
	Fatalf(ctx, err, printed)
	assert.Empty(t, entries(ctx))
	expected := buffered + err.Error() + printed
	assert.Equal(t, expected, buf.String())

	if !exitCalled {
		t.Error("exit not called")
	}
}

func TestDebug(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(testFormat))

	// Debug logs are ignored by default.
	Debugf(ctx, ignored)
	assert.Empty(t, entries(ctx))
	assert.Empty(t, buf.String())

	// Debug logs are enabled after setting the WithDebug option.
	ctx = Context(ctx, WithDebug())
	Debugf(ctx, printed)
	assert.Empty(t, entries(ctx))
	assert.Equal(t, printed, buf.String())

	// Buffering is disabled in debug mode.
	Infof(ctx, printed)
	assert.Empty(t, entries(ctx))
	assert.Equal(t, printed+printed, buf.String())
}

func TestStructuredLogging(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(testFormat))

	// No key-value pair is logged by default.
	Infof(ctx, buffered)
	if len(entries(ctx)) != 1 {
		t.Fatalf("got %d buffered entries, want 1", len(entries(ctx)))
	}
	e := (entries(ctx))[0]
	if len(e.KeyVals) != 1 {
		t.Errorf("got %d keyvals, want 1", len(e.KeyVals))
	}

	// Key-value pairs are logged.
	Info(ctx, KV{"key1", "val1"}, KV{"key2", "val2"})
	if len(entries(ctx)) != 2 {
		t.Fatalf("got %d buffered entries, want 2", len(entries(ctx)))
	}
	e = (entries(ctx))[1]
	if len(e.KeyVals) != 2 {
		t.Errorf("got %d keyvals, want 2", len(e.KeyVals))
	}
	if e.KeyVals[0].K != "key1" || e.KeyVals[0].V != "val1" {
		t.Errorf("got keyval %q=%q, want key1=val1", e.KeyVals[0].K, e.KeyVals[0].V)
	}
	if e.KeyVals[1].K != "key2" || e.KeyVals[1].V != "val2" {
		t.Errorf("got keyval %q=%q, want key2=val2", e.KeyVals[1].K, e.KeyVals[1].V)
	}

	// Key-value pairs set in the log context are logged.
	ctx = With(ctx, KV{"key1", "val1"}, KV{"key2", "val2"})
	Info(ctx, KV{"msg", buffered})
	if len(entries(ctx)) != 3 {
		t.Fatalf("got %d buffered entries, want 3", len(entries(ctx)))
	}
	e = (entries(ctx))[2]
	if len(e.KeyVals) != 3 {
		t.Errorf("got %d keyvals, want 3", len(e.KeyVals))
	}
	if e.KeyVals[0].K != "key1" || e.KeyVals[0].V != "val1" {
		t.Errorf("got keyval %q=%q, want key1=val1", e.KeyVals[0].K, e.KeyVals[0].V)
	}
	if e.KeyVals[1].K != "key2" || e.KeyVals[1].V != "val2" {
		t.Errorf("got keyval %q=%q, want key2=val2", e.KeyVals[1].K, e.KeyVals[1].V)
	}
	if e.KeyVals[2].K != "msg" || e.KeyVals[2].V != buffered {
		t.Errorf("got keyval %q=%q, want msg=buffered", e.KeyVals[2].K, e.KeyVals[2].V)
	}

	// Key-value pairs set in the log context prefix logged key/value pairs.
	Info(ctx, KV{"key3", "val3"}, KV{"key4", "val4"})
	if len(entries(ctx)) != 4 {
		t.Fatalf("got %d buffered entries, want 4", len(entries(ctx)))
	}
	e = (entries(ctx))[3]
	if len(e.KeyVals) != 4 {
		t.Errorf("got %d keyvals, want 4", len(e.KeyVals))
	}
	for i := 0; i < 4; i++ {
		suffix := fmt.Sprintf("%d", i+1)
		if e.KeyVals[i].K != "key"+suffix || e.KeyVals[i].V != "val"+suffix {
			t.Errorf("got keyval %q=%q, want key"+suffix+"=val"+suffix, e.KeyVals[i].K, e.KeyVals[i].V)
		}
	}

	// Key-value pairs set in the log context are logged in order they are set.
	ctx = With(ctx, KV{"key3", "val3"}, KV{"key4", "val4"})
	Info(ctx, KV{"msg", buffered})
	if len(entries(ctx)) != 5 {
		t.Fatalf("got %d buffered entries, want 5", len(entries(ctx)))
	}
	e = (entries(ctx))[4]
	if len(e.KeyVals) != 5 {
		t.Errorf("got %d keyvals, want 5", len(e.KeyVals))
	}
	for i := 0; i < 4; i++ {
		suffix := fmt.Sprintf("%d", i+1)
		if e.KeyVals[i].K != "key"+suffix || e.KeyVals[i].V != "val"+suffix {
			t.Errorf("got keyval %q=%q, want key"+suffix+"=val"+suffix, e.KeyVals[i].K, e.KeyVals[i].V)
		}
	}
	if e.KeyVals[4].K != "msg" || e.KeyVals[4].V != buffered {
		t.Errorf("got keyval %q=%q, want msg=buffered", e.KeyVals[4].K, e.KeyVals[4].V)
	}
}

func TestDynamicKeyVals(t *testing.T) {
	var buf bytes.Buffer
	kvfunc := func(ctx context.Context) []KV {
		return []KV{{"key1", "val1"}, {"key2", "val2"}}
	}
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(testFormat), WithFunc(kvfunc))
	Infof(ctx, buffered)
	if len(entries(ctx)) != 1 {
		t.Fatalf("got %d buffered entries, want 1", len(entries(ctx)))
	}
	e := (entries(ctx))[0]
	if len(e.KeyVals) != 3 {
		t.Errorf("got %d keyvals, want 3", len(e.KeyVals))
	}
	if e.KeyVals[0].K != "msg" || e.KeyVals[0].V != "buffered" {
		t.Errorf("got keyval %q=%q, want msg=buffered", e.KeyVals[0].K, e.KeyVals[0].V)
	}
	if e.KeyVals[1].K != "key1" || e.KeyVals[1].V != "val1" {
		t.Errorf("got keyval %q=%q, want key1=val1", e.KeyVals[1].K, e.KeyVals[1].V)
	}
	if e.KeyVals[2].K != "key2" || e.KeyVals[2].V != "val2" {
		t.Errorf("got keyval %q=%q, want key2=val2", e.KeyVals[2].K, e.KeyVals[2].V)
	}
}

func TestChaining(t *testing.T) {
	ctx1 := Context(context.Background())
	ctx2 := With(ctx1, KV{"key1", "val1"})
	Info(ctx1, KV{"msg", "msg1"})
	Info(ctx2, KV{"msg", "msg2"})

	if len(entries(ctx1)) != 1 {
		t.Fatalf("got %d buffered entries, want 1", len(entries(ctx1)))
	}
	e := (entries(ctx1))[0]
	if len(e.KeyVals) != 1 {
		t.Errorf("got %d keyvals, want 1", len(e.KeyVals))
	}

	if len(entries(ctx2)) != 1 {
		t.Fatalf("got %d buffered entries, want 1", len(entries(ctx2)))
	}
	e = (entries(ctx2))[0]
	if len(e.KeyVals) != 2 {
		t.Errorf("got %d keyvals, want 2", len(e.KeyVals))
	}
	if e.KeyVals[0].K != "key1" || e.KeyVals[0].V != "val1" {
		t.Errorf("got keyval %q=%q, want key1=val1", e.KeyVals[0].K, e.KeyVals[0].V)
	}
}

func TestNoLogging(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("panic: %v", err)
		}
	}()
	ctx := context.Background()
	Debugf(ctx, "")
	Printf(ctx, "")
	Infof(ctx, "")
	Errorf(ctx, nil, "")
	With(ctx, KV{"key", "val"})
	FlushAndDisableBuffering(ctx)
}

func TestMaxSize(t *testing.T) {
	var (
		txt            = "|txt|"
		maxsize        = len(txt)
		maxtruncated   = len(txt) + len(truncationSuffix)
		msg            = KV{"msg", "|txt|"}
		toolong        = KV{"msg", txt + "b"}
		toomany        = make([]string, maxsize+1)
		toomanytoolong = make([]KV, maxsize+1)
		toomanyi       = make([]KV, maxsize+1)
	)
	for i := 0; i < maxsize+1; i += 1 {
		idx := strconv.Itoa(i)
		toomany[i] = txt
		toomanytoolong[i] = KV{"key" + idx, txt + "b"}
		toomanyi[i] = KV{"key" + idx, interface{}(txt + "b")}
	}
	cases := []struct {
		name     string
		keyvals  []KV
		expected int
	}{
		{"short message", []KV{msg}, len(txt)},
		{"long message", []KV{toolong}, maxtruncated},
		{"too many elements in value", []KV{{"key", toomany}}, maxtruncated},
		{"too many too long elements in value", toomanytoolong, maxtruncated*maxsize + len(",log:"+truncationSuffix)},
		{"too many too long elements in []interface{} value", toomanyi, maxtruncated*maxsize + len(",log:"+truncationSuffix)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Append values
			format := func(e *Entry) []byte {
				var vals string
				for i := 0; i < len(e.KeyVals); i++ {
					if sv, ok := e.KeyVals[i].V.([]string); ok {
						vals += strings.Join(sv, "")
					} else if sv, ok := e.KeyVals[i].V.([]interface{}); ok {
						strs := make([]string, len(sv))
						for j := range sv {
							strs[j] = sv[j].(string)
						}
						vals += strings.Join(strs, "")
					} else {
						vals += e.KeyVals[i].V.(string)
					}
				}
				return []byte(vals)
			}

			ctx := Context(context.Background(), WithOutput(&buf), WithMaxSize(maxsize), WithFormat(format))
			Print(ctx, kvList(c.keyvals))

			if buf.Len() != c.expected {
				t.Errorf("got %d (%q), want %d", buf.Len(), buf.String(), c.expected)
			}
		})
	}

	t.Run("example result", func(t *testing.T) {
		now, epoc := timeNow, epoch
		timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
		epoch = timeNow()
		defer func() { timeNow = now; epoch = epoc }()

		var buf bytes.Buffer
		ctx := Context(context.Background(), WithOutput(&buf), WithMaxSize(maxsize), WithFormat(FormatText))
		Print(ctx, KV{"truncated", "it is too long"})

		want := "time=2022-01-09T20:29:45Z level=info truncated=\"it is ... <clue/log.truncated>\"\n"
		if got := buf.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestWithMultipleDerivedLoggers(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatText))
	FlushAndDisableBuffering(ctx)

	Info(ctx, KV{"msg", "root"})

	ctx1 := With(ctx, KV{"key1", "value1"})
	Info(ctx1, KV{"msg", "first"})

	ctx2 := With(ctx, KV{"key2", "value2"})
	Info(ctx2, KV{"msg", "second"})
	Info(ctx1, KV{"msg", "third"})

	ctx3 := With(ctx1, KV{"key3", "value3"})
	Info(ctx3, KV{"msg", "fourth"})
	Info(ctx2, KV{"msg", "fifth"})
	Info(ctx1, KV{"msg", "sixth"})
	Info(ctx, KV{"msg", "seventh"})

	logs := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, logs, 8)
	assert.Equal(t, "time=2022-02-22T17:00:00Z level=info msg=root", logs[0])
	assert.Equal(t, "time=2022-02-22T17:00:00Z level=info key1=value1 msg=first", logs[1])
	assert.Equal(t, "time=2022-02-22T17:00:00Z level=info key2=value2 msg=second", logs[2])
	assert.Equal(t, "time=2022-02-22T17:00:00Z level=info key1=value1 msg=third", logs[3])
	assert.Equal(t, "time=2022-02-22T17:00:00Z level=info key1=value1 key3=value3 msg=fourth", logs[4])
	assert.Equal(t, "time=2022-02-22T17:00:00Z level=info key2=value2 msg=fifth", logs[5])
	assert.Equal(t, "time=2022-02-22T17:00:00Z level=info key1=value1 msg=sixth", logs[6])
	assert.Equal(t, "time=2022-02-22T17:00:00Z level=info msg=seventh", logs[7])
}

func TestConcurrentLog(t *testing.T) {
	// run with -race to detect data races
	var buf Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(FormatText))
	FlushAndDisableBuffering(ctx)

	for i := 0; i < 10; i++ {
		go func() {
			Info(ctx, KV{"msg", "root"})
			ctx1 := With(ctx, KV{"key1", "value1"})
			Info(ctx1, KV{"msg", "first"})
		}()
	}
}

func testFormat(e *Entry) []byte {
	var buf bytes.Buffer
	for _, kv := range e.KeyVals {
		fmt.Fprintf(&buf, "%s", kv.V)
	}
	return buf.Bytes()
}

func entries(ctx context.Context) []*Entry {
	l := ctx.Value(ctxLogger).(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.entries
}
