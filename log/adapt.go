package log

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/smithy-go/logging"
	"github.com/go-logr/logr"
	"goa.design/goa/v3/middleware"
)

type (
	// StdLogger implements an interface compatible with the stdlib log package.
	StdLogger struct {
		ctx context.Context
	}

	// AWSLogger returns an AWS SDK compatible logger.
	AWSLogger struct {
		context.Context
	}

	// LogrSink returns a logr LogSink compatible logger.
	LogrSink struct {
		lock sync.Mutex
		name string
		context.Context
	}

	// goaLogger is a Goa middleware compatible logger.
	goaLogger struct {
		context.Context
	}
)

// NameKey is the key used to log the name of the logger.
const NameKey = "log"

// AsGoaMiddlewareLogger creates a Goa middleware compatible logger that can be used when
// configuring Goa HTTP or gRPC servers.
//
// Usage:
//
//	// HTTP server:
//	import goahttp "goa.design/goa/v3/http"
//	import httpmdlwr "goa.design/goa/v3/http/middleware"
//	...
//	mux := goahttp.NewMuxer()
//	handler := httpmdlwr.LogContext(log.AsGoaMiddlewareLogger)(mux)
//
//	// gRPC server:
//	import "google.golang.org/grpc"
//	import grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
//	import grpcmdlwr "goa.design/goa/v3/grpc/middleware"
//	...
//	srv := grpc.NewServer(
//	    grpcmiddleware.WithUnaryServerChain(grpcmdlwr.UnaryServerLogContext(log.AsGoaMiddlewareLogger)),
//	)
func AsGoaMiddlewareLogger(ctx context.Context) middleware.Logger {
	return goaLogger{ctx}
}

// AsStdLogger adapts a Goa logger to a stdlib compatible logger.
func AsStdLogger(ctx context.Context) *StdLogger {
	return &StdLogger{ctx}
}

// AsAWSLogger returns an AWS SDK compatible logger.
//
// Usage:
//
//	import "github.com/aws/aws-sdk-go-v2/config"
//	import "goa.design/clue/log"
//	import "goa.design/clue/trace"
//
//	ctx := log.Context(context.Background())
//	tracedClient := &http.Client{Transport: trace.Client(ctx, http.DefaultTransport)}
//	cfg, err := config.LoadDefaultConfig(ctx,
//	    config.WithHTTPClient(tracedClient),
//	    config.WithLogger(log.AsAWSLogger(ctx)))
func AsAWSLogger(ctx context.Context) *AWSLogger {
	return &AWSLogger{ctx}
}

// ToLogrSink returns a logr.LogSink.
//
// Usage:
//
//	import "goa.design/clue/log"
//
//	ctx := log.Context(context.Background())
//	sink := log.ToLogrSink(ctx)
//	logger := logr.New(sink)
func ToLogrSink(ctx context.Context) *LogrSink {
	return &LogrSink{Context: ctx}
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *StdLogger) Fatal(v ...interface{}) {
	l.Print(v...)
	osExit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *StdLogger) Fatalf(format string, v ...interface{}) {
	l.Printf(format, v...)
	osExit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *StdLogger) Fatalln(v ...interface{}) {
	l.Println(v...)
	osExit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *StdLogger) Panic(v ...interface{}) {
	l.Print(v...)
	panic(fmt.Sprint(v...))
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *StdLogger) Panicf(format string, v ...interface{}) {
	l.Printf(format, v...)
	panic(fmt.Sprintf(format, v...))
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *StdLogger) Panicln(v ...interface{}) {
	l.Println(v...)
	panic(fmt.Sprintln(v...))
}

// Print print to the logger. Arguments are handled in the manner of fmt.Print.
func (l *StdLogger) Print(v ...interface{}) {
	Printf(l.ctx, "%s", fmt.Sprint(v...))
}

// Printf prints to the logger. Arguments are handled in the manner of fmt.Printf.
func (l *StdLogger) Printf(format string, v ...interface{}) {
	Printf(l.ctx, format, v...)
}

// Println prints to the logger. Arguments are handled in the manner of fmt.Println.
func (l *StdLogger) Println(v ...interface{}) {
	Printf(l.ctx, "%s", fmt.Sprintln(v...))
}

func (l *AWSLogger) Logf(classification logging.Classification, format string, v ...any) {
	fn := Infof
	if classification == logging.Debug {
		fn = Debugf
	}
	fn(l, format, v...)
}

func (l *AWSLogger) WithContext(ctx context.Context) logging.Logger {
	l.Context = WithContext(ctx, l)
	return l
}

func (l *LogrSink) Init(info logr.RuntimeInfo) {}

func (l *LogrSink) Enabled(level int) bool {
	return true
}

func (l *LogrSink) Info(level int, msg string, keysAndValues ...interface{}) {
	kvs := make([]KV, len(keysAndValues)/2+1)
	kvs[0] = KV{K: "msg", V: msg}
	for i := 0; i < len(keysAndValues); i += 2 {
		kvs[i/2+1] = KV{K: fmt.Sprint(keysAndValues[i]), V: keysAndValues[i+1]}
	}
	if level == 0 {
		Info(l, kvList(kvs))
	} else {
		Debug(l, kvList(kvs))
	}
}

func (l *LogrSink) Error(err error, msg string, keysAndValues ...interface{}) {
	kvs := make([]KV, len(keysAndValues)/2+1)
	kvs[0] = KV{K: "msg", V: msg}
	for i := 0; i < len(keysAndValues); i += 2 {
		kvs[i/2+1] = KV{K: fmt.Sprint(keysAndValues[i]), V: keysAndValues[i+1]}
	}
	Error(l, err, kvList(kvs))
}

func (l *LogrSink) WithValues(keysAndValues ...any) logr.LogSink {
	kvs := make([]KV, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		kvs[i/2] = KV{K: fmt.Sprint(keysAndValues[i]), V: keysAndValues[i+1]}
	}
	return &LogrSink{Context: With(l, kvList(kvs))}
}

func (l *LogrSink) WithName(name string) logr.LogSink {
	l.lock.Lock()
	defer l.lock.Unlock()
	cur := l.name
	if cur != "" {
		cur += "/"
	}
	l.name = cur + name
	return &LogrSink{Context: With(l, KV{NameKey, l.name}), name: l.name}
}

// Log creates a log entry using a sequence of key/value pairs.
func (l goaLogger) Log(keyvals ...interface{}) error {
	n := (len(keyvals) + 1) / 2
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "MISSING")
	}
	kvs := make([]KV, n)
	for i := 0; i < n; i++ {
		k, v := keyvals[2*i], keyvals[2*i+1]
		kvs[i] = KV{K: fmt.Sprint(k), V: v}
	}
	Print(l, kvList(kvs))
	return nil
}
