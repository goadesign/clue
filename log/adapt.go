package log

import (
	"context"
	"fmt"

	"goa.design/goa/v3/middleware"
)

type (
	// goaLogger is a Goa middleware compatible logger.
	goaLogger struct {
		context.Context
	}
)

// Adapt creates a Goa middleware compatible logger that can be used when
// configuring Goa HTTP or gRPC servers.
//
// Usage:
//
//    ctx := log.Context(context.Background())
//    logger := log.Adapt(ctx)
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
func Adapt(ctx context.Context) middleware.Logger {
	return goaLogger{ctx}
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
	Print(l, kvs...)
	return nil
}
