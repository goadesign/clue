package log

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestAsGoaMiddlwareLogger(t *testing.T) {
	restore := timeNow
	timeNow = func() time.Time { return time.Date(2022, time.January, 9, 20, 29, 45, 0, time.UTC) }
	defer func() { timeNow = restore }()

	var buf bytes.Buffer
	ctx := Context(context.Background(), WithOutput(&buf))
	logger := AsGoaMiddlewareLogger(ctx)
	logger.Log("msg", "hello world")
	want := "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
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
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}

	buf.Reset()
	logger.Println("hello world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\\n\"\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}

	buf.Reset()
	logger.Printf("hello %s", "world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}

	func() {
		buf.Reset()
		var msg string
		defer func() { msg = recover().(string) }()
		logger.Panic("hello world")
		want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
		if buf.String() != want {
			t.Errorf("got %q, want %q", buf.String(), want)
		}
		if msg != "hello world" {
			t.Errorf("got %q, want %q", msg, "hello world")
		}
	}()

	func() {
		buf.Reset()
		var msg string
		defer func() { msg = recover().(string) }()
		logger.Panicf("hello %s", "world")
		want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
		if buf.String() != want {
			t.Errorf("got %q, want %q", buf.String(), want)
		}
		if msg != "hello world" {
			t.Errorf("got %q, want %q", msg, "hello world")
		}
	}()

	func() {
		buf.Reset()
		var msg string
		defer func() { msg = recover().(string) }()
		logger.Panicln("hello world")
		want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\\n\"\n"
		if buf.String() != want {
			t.Errorf("got %q, want %q", buf.String(), want)
		}
		if msg != "hello world\n" {
			t.Errorf("got %q, want %q", msg, "hello world\n")
		}
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
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
	if exited != 1 {
		t.Errorf("got %d, want %d", exited, 1)
	}

	exited = 0
	buf.Reset()
	logger.Fatalf("hello %s", "world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\"\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
	if exited != 1 {
		t.Errorf("got %d, want %d", exited, 1)
	}

	exited = 0
	buf.Reset()
	logger.Fatalln("hello world")
	want = "time=2022-01-09T20:29:45Z level=info msg=\"hello world\\n\"\n"
	if buf.String() != want {
		t.Errorf("got %q, want %q", buf.String(), want)
	}
	if exited != 1 {
		t.Errorf("got %d, want %d", exited, 1)
	}
}
