package log

import (
	"context"
	"fmt"

	"goa.design/goa/v3/middleware"
)

type (
	// StdLogger implements an interface compatible with the stdlib log package.
	StdLogger struct {
		ctx context.Context
	}

	// goaLogger is a Goa middleware compatible logger.
	goaLogger struct {
		context.Context
	}
)

// AsGoaMiddlewareLogger creates a Goa middleware compatible logger that can be used when
// configuring Goa HTTP or gRPC servers.
//
// Usage:
//
//    ctx := log.Context(context.Background())
//    logger := log.AsGoaMiddlewareLogger(ctx)
//
//    // HTTP server:
//    import goahttp "goa.design/goa/v3/http"
//    import httpmdlwr "goa.design/goa/v3/http/middleware"
//    ...
//    mux := goahttp.NewMuxer()
//    handler := httpmdlwr.Log(logger)(mux)
//
//    // gRPC server:
//    import "google.golang.org/grpc"
//    import grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
//    import grpcmdlwr "goa.design/goa/v3/grpc/middleware"
//    ...
//    srv := grpc.NewServer(
//        grpcmiddleware.WithUnaryServerChain(grpcmdlwr.UnaryServerLog(logger))
//    )
func AsGoaMiddlewareLogger(ctx context.Context) middleware.Logger {
	return goaLogger{ctx}
}

// AsStdLogger adapts a Goa logger to a stdlib compatible logger.
func AsStdLogger(ctx context.Context) *StdLogger {
	return &StdLogger{ctx}
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
