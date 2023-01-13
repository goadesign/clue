package debug

import (
	"context"
	"net/http"
	"net/http/pprof"

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
	// wantDebugEnabled is true if debug logs should be enabled.
	wantDebugEnabled bool

	// debugEnabled is true if debug logs are currently enabled.
	debugEnabled bool

	// pprofEnabled is true if pprof handlers are enabled.
	pprofEnabled bool
)

// MountDebugLogEnabler mounts an endpoint under the given prefix and returns a
// HTTP middleware that manages debug logs. The endpoint accepts a single query
// parameter "enable". If the parameter is set to "on" then debug logs are
// enabled for requests made to handlers returned by the middleware. If the
// parameter is set to "off" then debug logs are disabled.
func MountDebugLogEnabler(prefix string, mux Muxer) func(http.Handler) http.Handler {
	mux.Handle(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("enable") {
		case "on":
			wantDebugEnabled = true
		case "off":
			wantDebugEnabled = false
		}
	}))
	return func(next http.Handler) http.Handler {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wantDebugEnabled && !debugEnabled {
				r = r.WithContext(log.Context(r.Context(), log.WithDebug()))
				debugEnabled = true
			} else if !wantDebugEnabled && debugEnabled {
				r = r.WithContext(log.Context(r.Context(), log.WithNoDebug()))
				debugEnabled = false
			}
			next.ServeHTTP(w, r)
		})
		return handler
	}
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
// Note: do not call this function on production servers accessible to the
// public!  It exposes sensitive information about the server.
func MountPprofHandlers(mux Muxer) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

// LogPayloads returns a Goa endpoint middleware that logs request payloads and
// response results using debug log entries.
//
// Note: this middleware marshals the request and response data using the
// standard JSON marshaller. It only marshals if debug logs are enabled.
func LogPayloads(opts ...LogPayloadsOption) func(goa.Endpoint) goa.Endpoint {
	return func(next goa.Endpoint) goa.Endpoint {
		options := defaultOptions()
		for _, opt := range opts {
			if opt != nil {
				opt(options)
			}
		}
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if !log.DebugEnabled(ctx) {
				return next(ctx, req)
			}
			reqs := options.format(ctx, req)
			if len(reqs) > options.maxsize {
				reqs = reqs[:options.maxsize]
			}
			log.Debug(ctx, log.KV{K: "payload", V: reqs})
			res, err := next(ctx, req)
			if err != nil {
				return nil, err
			}
			ress := options.format(ctx, res)
			if len(ress) > options.maxsize {
				ress = ress[:options.maxsize]
			}
			log.Debug(ctx, log.KV{K: "result", V: ress})
			return res, nil
		}
	}
}
