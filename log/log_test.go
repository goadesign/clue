package log

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

const (
	buffered = "buffered"
	printed  = "printed"
	ignored  = "ignored"
)

func TestKeyValParse(t *testing.T) {
	cases := []struct {
		name         string
		keyvals      []interface{}
		expectedKeys []string
		expectedVals []interface{}
	}{
		{"empty", []interface{}{}, []string{}, []interface{}{}},
		{"one", []interface{}{"key", "val"}, []string{"key"}, []interface{}{"val"}},
		{"invalid key", []interface{}{0, "val"}, []string{"<INVALID>"}, []interface{}{"val"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			kv := KeyVals(c.keyvals)
			keys, vals := kv.Parse()
			if len(keys) != len(c.expectedKeys) {
				t.Fatalf("got %d keys, want %d", len(keys), len(c.expectedKeys))
			}
			if len(vals) != len(c.expectedVals) {
				t.Fatalf("got %d vals, want %d", len(vals), len(c.expectedVals))
			}
			for i, k := range keys {
				if k != c.expectedKeys[i] {
					t.Errorf("got key %q, want %q", k, c.expectedKeys[i])
				}
			}
			for i, v := range vals {
				if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", c.expectedVals[i]) {
					t.Errorf("got val %v, want %v", v, c.expectedVals[i])
				}
			}
		})
	}
}

func TestSeverity(t *testing.T) {
	var buf bytes.Buffer
	printSev := func(e *Entry) []byte {
		return []byte(e.Severity.String() + ":" + e.Severity.Code() + ":" + e.Severity.Color() + " ")
	}
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(printSev), WithDebug())
	Debug(ctx, "")
	Info(ctx, "")
	Error(ctx, "")
	want := "DEBUG:DEBG:\033[37m INFO:INFO:\033[34m ERROR:ERRO:\033[1;31m "
	if got := buf.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if Severity(0).String() != "<INVALID>" {
		t.Errorf("got %q, want %q", Severity(0).String(), "<INVALID>")
	}
	if Severity(0).Code() != "<INVALID>" {
		t.Errorf("got %q, want %q", Severity(0).Code(), "<INVALID>")
	}
	if Severity(0).Color() != "" {
		t.Errorf("got %q, want empty", Severity(0).Color())
	}
}

func TestBuffering(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(debugFormat))

	// Buffering is enabled by default.
	Info(ctx, buffered)
	if len(entries(ctx)) != 1 {
		t.Errorf("got %d buffered entries, want 1", len(entries(ctx)))
	} else {
		e := entries(ctx)[0]
		if e.Message != buffered {
			t.Errorf("got buffered entry message %q, want %q", e.Message, buffered)
		}
	}

	// Print does not buffer.
	Print(ctx, printed)
	if buf.String() != printed {
		t.Errorf("got printed message %q, want %q", buf.String(), printed)
	}

	// Flush flushes the buffer.
	Flush(ctx)
	if len(entries(ctx)) != 0 {
		t.Errorf("got %d buffered entries, want 0", len(entries(ctx)))
	}
	if buf.String() != printed+buffered {
		t.Errorf("got printed message %q, want %q", buf.String(), printed+buffered)
	}

	// Buffering is disabled after flush.
	Info(ctx, printed)
	if len(entries(ctx)) != 0 {
		t.Errorf("got %d buffered entries, want 0", len(entries(ctx)))
	}
	if buf.String() != printed+buffered+printed {
		t.Errorf("got printed message %q, want %q", buf.String(), printed+buffered+printed)
	}

	// Flush is idempotent.
	Flush(ctx)
	Info(ctx, printed)
	if len(entries(ctx)) != 0 {
		t.Errorf("got %d buffered entries, want 0", len(entries(ctx)))
	}
	if buf.String() != printed+buffered+printed+printed {
		t.Errorf("got printed message %q, want %q", buf.String(), printed+buffered+printed+printed)
	}
}

func TestBufferingWithError(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(debugFormat))

	// Error flushes the buffer.
	Info(ctx, buffered)
	Error(ctx, printed)
	if len(entries(ctx)) != 0 {
		t.Errorf("got %d buffered entries, want 0", len(entries(ctx)))
	}
	if buf.String() != buffered+printed {
		t.Errorf("got printed message %q, want %q", buf.String(), buffered+printed)
	}

	// Buffering is disabled after error.
	Info(ctx, printed)
	if len(entries(ctx)) != 0 {
		t.Errorf("got %d buffered entries, want 0", len(entries(ctx)))
	}
	if buf.String() != buffered+printed+printed {
		t.Errorf("got printed message %q, want %q", buf.String(), buffered+printed+printed)
	}
}

func TestDebug(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(debugFormat))

	// Debug logs are ignored by default.
	Debug(ctx, ignored)
	if len(entries(ctx)) != 0 {
		t.Errorf("got %d buffered entries, want 0", len(entries(ctx)))
	}
	if buf.String() != "" {
		t.Errorf("got printed message %q, want empty", buf.String())
	}

	// Debug logs are enabled after setting the WithDebug option.
	ctx = Context(ctx, WithDebug())
	Debug(ctx, printed)
	if len(entries(ctx)) != 0 {
		t.Errorf("got %d buffered entries, want 0", len(entries(ctx)))
	}
	if buf.String() != printed {
		t.Errorf("got printed message %q, want %q", buf.String(), printed)
	}

	// Buffering is disabled in debug mode.
	Info(ctx, printed)
	if len(entries(ctx)) != 0 {
		t.Errorf("got %d buffered entries, want 0", len(entries(ctx)))
	}
	if buf.String() != printed+printed {
		t.Errorf("got printed message %q, want %q", buf.String(), printed+printed)
	}
}

func TestStructuredLogging(t *testing.T) {
	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf), WithFormat(debugFormat))

	// No key-value pair is logged by default.
	Info(ctx, buffered)
	if len(entries(ctx)) != 1 {
		t.Fatalf("got %d buffered entries, want 1", len(entries(ctx)))
	}
	e := (entries(ctx))[0]
	if len(e.KeyVals) != 0 {
		t.Errorf("got %d keyvals, want 0", len(e.KeyVals))
	}

	// Key-value pairs are logged.
	Info(ctx, buffered, "key1", "val1", "key2", "val2")
	if len(entries(ctx)) != 2 {
		t.Fatalf("got %d buffered entries, want 2", len(entries(ctx)))
	}
	e = (entries(ctx))[1]
	if len(e.KeyVals) != 4 {
		t.Errorf("got %d keyvals, want 4", len(e.KeyVals))
	}
	keys, vals := e.KeyVals.Parse()
	if keys[0] != "key1" || vals[0] != "val1" {
		t.Errorf("got keyval %q=%q, want key1=val1", keys[0], vals[0])
	}
	if keys[1] != "key2" || vals[1] != "val2" {
		t.Errorf("got keyval %q=%q, want key2=val2", keys[1], vals[1])
	}

	// log does not panic when an odd number of arguments is given to Info.
	Info(ctx, buffered, "key1")
	if len(entries(ctx)) != 3 {
		t.Fatalf("got %d buffered entries, want 3", len(entries(ctx)))
	}
	e = (entries(ctx))[2]
	if len(e.KeyVals) != 2 {
		t.Errorf("got %d keyvals, want 2", len(e.KeyVals))
	}
	keys, vals = e.KeyVals.Parse()
	if keys[0] != "key1" || vals[0] != nil {
		t.Errorf("got keyval %q=%q, want key1=", keys[0], vals[0])
	}

	// Key-value pairs set in the log context are logged.
	ctx = With(ctx, "key1", "val1", "key2", "val2")
	Info(ctx, buffered)
	if len(entries(ctx)) != 4 {
		t.Fatalf("got %d buffered entries, want 4", len(entries(ctx)))
	}
	e = (entries(ctx))[3]
	if len(e.KeyVals) != 4 {
		t.Errorf("got %d keyvals, want 4", len(e.KeyVals))
	}
	keys, vals = e.KeyVals.Parse()
	if keys[0] != "key1" || vals[0] != "val1" {
		t.Errorf("got keyval %q=%q, want key1=val1", keys[0], vals[0])
	}
	if keys[1] != "key2" || vals[1] != "val2" {
		t.Errorf("got keyval %q=%q, want key2=val2", keys[1], vals[1])
	}

	// Key-value pairs set in the log context prefix logged key/value pairs.
	Info(ctx, buffered, "key3", "val3", "key4", "val4")
	if len(entries(ctx)) != 5 {
		t.Fatalf("got %d buffered entries, want 5", len(entries(ctx)))
	}
	e = (entries(ctx))[4]
	if len(e.KeyVals) != 8 {
		t.Errorf("got %d keyvals, want 8", len(e.KeyVals))
	}
	keys, vals = e.KeyVals.Parse()
	for i := 0; i < 4; i++ {
		suffix := fmt.Sprintf("%d", i+1)
		if keys[i] != "key"+suffix || vals[i] != "val"+suffix {
			t.Errorf("got keyval %q=%q, want key"+suffix+"=val"+suffix, keys[i], vals[i])
		}
	}

	// Key-value pairs set in the log context are logged in order they are set.
	ctx = With(ctx, "key3", "val3", "key4", "val4")
	Info(ctx, buffered)
	if len(entries(ctx)) != 6 {
		t.Fatalf("got %d buffered entries, want 6", len(entries(ctx)))
	}
	e = (entries(ctx))[5]
	if len(e.KeyVals) != 8 {
		t.Errorf("got %d keyvals, want 8", len(e.KeyVals))
	}
	keys, vals = e.KeyVals.Parse()
	for i := 0; i < 4; i++ {
		suffix := fmt.Sprintf("%d", i+1)
		if keys[i] != "key"+suffix || vals[i] != "val"+suffix {
			t.Errorf("got keyval %q=%q, want key"+suffix+"=val"+suffix, keys[i], vals[i])
		}
	}

	// log does not panic when an odd number of arguments is given to With.
	ctx = With(ctx, "key3")
	Info(ctx, buffered)
	if len(entries(ctx)) != 7 {
		t.Fatalf("got %d buffered entries, want 7", len(entries(ctx)))
	}
	e = (entries(ctx))[6]
	if len(e.KeyVals) != 10 {
		t.Errorf("got %d keyvals, want 10", len(e.KeyVals))
	}
	keys, vals = e.KeyVals.Parse()
	if len(keys) != 5 {
		t.Fatalf("got %d keys, want 5", len(keys))
	}
	if len(vals) != 5 {
		t.Fatalf("got %d vals, want 5", len(vals))
	}
	if keys[4] != "key3" || vals[4] != nil {
		t.Errorf("got keyval %q=%q, want key3=", keys[4], vals[4])
	}
}

func TestChaining(t *testing.T) {
	ctx1 := Context(context.Background())
	ctx2 := With(ctx1, "key1", "val1")
	Info(ctx1, "msg1")
	Info(ctx2, "msg2")

	if len(entries(ctx1)) != 1 {
		t.Fatalf("got %d buffered entries, want 1", len(entries(ctx1)))
	}
	e := (entries(ctx1))[0]
	if len(e.KeyVals) != 0 {
		t.Errorf("got %d keyvals, want 0", len(e.KeyVals))
	}

	if len(entries(ctx2)) != 1 {
		t.Fatalf("got %d buffered entries, want 1", len(entries(ctx2)))
	}
	e = (entries(ctx2))[0]
	if len(e.KeyVals) != 2 {
		t.Errorf("got %d keyvals, want 2", len(e.KeyVals))
	}
	keys, vals := e.KeyVals.Parse()
	if keys[0] != "key1" || vals[0] != "val1" {
		t.Errorf("got keyval %q=%q, want key1=val1", keys[0], vals[0])
	}
}

func TestNoLogging(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("panic: %v", err)
		}
	}()
	ctx := context.Background()
	Debug(ctx, "")
	Print(ctx, "")
	Info(ctx, "")
	Error(ctx, "")
	With(ctx, "key", "val")
	Flush(ctx)
}

func TestMaxSize(t *testing.T) {
	var (
		txt            = "|txt|"
		maxsize        = len(txt)
		keyval         = []interface{}{"key", txt}
		toolong        = []interface{}{"key", txt + "b"}
		toomany        = make([]string, maxsize+1)
		toomanytoolong = make([]string, maxsize+1)
		toomanyi       = make([]interface{}, maxsize+1)
	)
	for i := 0; i < maxsize+1; i += 1 {
		toomany[i] = txt
		toomanytoolong[i] = txt + "b"
		toomanyi[i] = txt + "b"
	}
	cases := []struct {
		name     string
		msg      string
		keyvals  []interface{}
		expected int
	}{
		{"short message", txt, nil, len(txt)},
		{"long message", txt + "a", nil, len(txt)},
		{"short message with short value", txt, keyval, 2 * len(txt)},
		{"long message with short value", txt + "a", keyval, 2 * len(txt)},
		{"short message with long value", txt, toolong, 2*len(txt) + len(truncationSuffix)},
		{"long message with long value", txt + "a", toolong, 2*len(txt) + len(truncationSuffix)},
		{"too many elements in value", "", []interface{}{"key", toomany}, maxsize + len(truncationSuffix)},
		{"too many too long elements in value", "", []interface{}{"key", toomanytoolong}, maxsize + len(truncationSuffix)},
		{"too many too long elements in []interface{} value", "", []interface{}{"key", toomanyi}, maxsize + len(truncationSuffix)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Append message and values
			format := func(e *Entry) []byte {
				var vals string
				for i := 1; i < len(e.KeyVals); i += 2 {
					if sv, ok := e.KeyVals[i].([]string); ok {
						vals += strings.Join(sv, "")
					} else if sv, ok := e.KeyVals[i].([]interface{}); ok {
						strs := make([]string, len(sv))
						for j := range sv {
							strs[j] = sv[j].(string)
						}
						vals += strings.Join(strs, "")
					} else {
						vals += e.KeyVals[i].(string)
					}
				}
				return []byte(e.Message + vals)
			}

			ctx := Context(context.Background(), WithOutput(&buf), WithMaxSize(maxsize), WithFormat(format))
			Print(ctx, c.msg, c.keyvals...)

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
		Print(ctx, "example", "truncated", "it is too long")

		want := "INFO[2022-01-09T20:29:45Z] examp truncated=it is ... <clue/log.truncated>\n"
		if got := buf.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func debugFormat(e *Entry) []byte {
	return []byte(e.Message)
}

func entries(ctx context.Context) []*Entry {
	l := ctx.Value(ctxLogger).(*logger)
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.entries
}
