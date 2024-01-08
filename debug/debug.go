package debug

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strings"

	goa "goa.design/goa/v3/pkg"

	"goa.design/clue/log"
)

// Muxer is the HTTP mux interface used by the debug package.
type Muxer interface {
	http.Handler
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

var (
	// debugLogs is true if debug logs should be enabled.
	debugLogs bool
)

// MountDebugLogEnabler mounts an endpoint under "/debug" that manages the
// status of debug logs. The endpoint accepts a single query parameter
// "debug-logs". If the parameter is set to "on" then debug logs are enabled. If
// the parameter is set to "off" then debug logs are disabled. In all other
// cases the endpoint returns the current debug logs status. The path, query
// parameter name and values can be changed using the WithPath, WithQuery,
// WithOnValue and WithOffValue options.
//
// Note: the endpoint merely controls the status of debug logs. It does not
// actually configure the current logger. The logger is configured by the
// middleware returned by the HTTP function or by the gRPC interceptors returned
// by the UnaryServerInterceptor and StreamServerInterceptor functions which
// should be used in conjunction with the MountDebugLogEnabler function.
func MountDebugLogEnabler(mux Muxer, opts ...DebugLogEnablerOption) {
	o := defaultDebugLogEnablerOptions()
	for _, opt := range opts {
		opt(o)
	}
	if !strings.HasPrefix(o.path, "/") {
		o.path = "/" + o.path
	}
	mux.Handle(o.path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if q := r.URL.Query().Get(o.query); q == o.onval {
			debugLogs = true
		} else if q == o.offval {
			debugLogs = false
		}
		if debugLogs {
			w.Write([]byte(fmt.Sprintf(`{"%s":"%s"}`, o.query, o.onval))) // nolint: errcheck
		} else {
			w.Write([]byte(fmt.Sprintf(`{"%s":"%s"}`, o.query, o.offval))) // nolint: errcheck
		}
	}))
}

// MountPprofHandlers mounts pprof handlers under /debug/pprof/. The list of
// mounted handlers is:
//
//	/debug/pprof/
//	/debug/pprof/allocs
//	/debug/pprof/block
//	/debug/pprof/cmdline
//	/debug/pprof/goroutine
//	/debug/pprof/heap
//	/debug/pprof/mutex
//	/debug/pprof/profile
//	/debug/pprof/symbol
//	/debug/pprof/threadcreate
//	/debug/pprof/trace
//
// See the pprof package documentation for more information.
//
// The path prefix ("/debug/pprof/") can be changed using WithPprofPrefix.
// Note: do not call this function on production servers accessible to the
// public!  It exposes sensitive information about the server.
func MountPprofHandlers(mux Muxer, opts ...PprofOption) {
	o := defaultPprofOptions()
	for _, opt := range opts {
		opt(o)
	}
	if !strings.HasPrefix(o.prefix, "/") {
		o.prefix = "/" + o.prefix
	}
	if !strings.HasSuffix(o.prefix, "/") {
		o.prefix = o.prefix + "/"
	}
	mux.HandleFunc(o.prefix, pprof.Index)
	mux.HandleFunc(o.prefix+"cmdline", pprof.Cmdline)
	mux.HandleFunc(o.prefix+"profile", pprof.Profile)
	mux.HandleFunc(o.prefix+"symbol", pprof.Symbol)
	mux.HandleFunc(o.prefix+"trace", pprof.Trace)
}

// LogPayloads returns a Goa endpoint middleware that logs request payloads and
// response results using debug log entries.
//
// Note: this middleware marshals the request and response data using the
// standard JSON marshaller. It only marshals if debug logs are enabled.
func LogPayloads(opts ...LogPayloadsOption) func(goa.Endpoint) goa.Endpoint {
	return func(next goa.Endpoint) goa.Endpoint {
		options := defaultLogPayloadsOptions()
		for _, opt := range opts {
			if opt != nil {
				opt(options)
			}
		}
		reqKey := "payload"
		resKey := "result"
		if options.client {
			reqKey = "client-" + reqKey
			resKey = "client-" + resKey
		}
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if !log.DebugEnabled(ctx) {
				return next(ctx, req)
			}
			reqs := options.format(ctx, req)
			if len(reqs) > options.maxsize {
				reqs = reqs[:options.maxsize]
			}
			log.Debug(ctx, log.KV{K: reqKey, V: reqs})
			res, err := next(ctx, req)
			if err != nil {
				return nil, err
			}
			ress := options.format(ctx, res)
			if len(ress) > options.maxsize {
				ress = ress[:options.maxsize]
			}
			log.Debug(ctx, log.KV{K: resKey, V: ress})
			return res, nil
		}
	}
}
