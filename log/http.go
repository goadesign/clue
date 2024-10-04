package log

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"

	goa "goa.design/goa/v3/pkg"
)

type (
	// HTTPLogOption is a function that applies a configuration option to log
	// HTTP middleware.
	HTTPLogOption func(*httpLogOptions)

	// HTTPClientLogOption is a function that applies a configuration option
	// to a HTTP client logger.
	HTTPClientLogOption func(*httpClientOptions)

	httpLogOptions struct {
		pathFilters           []*regexp.Regexp
		disableRequestLogging bool
		disableRequestID      bool
	}

	httpClientOptions struct {
		iserr      func(int) bool
		logErrBody bool
	}

	// client wraps an HTTP roundtripper and logs requests and responses.
	client struct {
		http.RoundTripper
		options *httpClientOptions
	}

	// responseCapture is a http.ResponseWriter which captures the response status
	// code and content length.
	responseCapture struct {
		http.ResponseWriter
		StatusCode    int
		ContentLength int
	}
)

// HTTP returns a HTTP middleware that performs two tasks:
//  1. Enriches the request context with the logger specified in logCtx.
//  2. Logs HTTP request details, except when WithDisableRequestLogging is set or
//     URL path matches a WithPathFilter regex.
//
// HTTP panics if logCtx was not created with Context.
func HTTP(logCtx context.Context, opts ...HTTPLogOption) func(http.Handler) http.Handler {
	MustContainLogger(logCtx)
	var options httpLogOptions
	for _, o := range opts {
		if o != nil {
			o(&options)
		}
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			for _, opt := range options.pathFilters {
				if opt.MatchString(req.URL.Path) {
					h.ServeHTTP(w, req)
					return
				}
			}
			ctx := WithContext(req.Context(), logCtx)
			if !options.disableRequestID {
				ctx = With(ctx, KV{RequestIDKey, shortID()})
			}
			if options.disableRequestLogging {
				h.ServeHTTP(w, req.WithContext(ctx))
				return
			}
			methKV := KV{K: HTTPMethodKey, V: req.Method}
			urlKV := KV{K: HTTPURLKey, V: req.URL.String()}
			fromKV := KV{K: HTTPFromKey, V: from(req)}
			Info(ctx, KV{K: MessageKey, V: "start"}, methKV, urlKV, fromKV)

			rw := &responseCapture{ResponseWriter: w}
			started := timeNow()
			h.ServeHTTP(rw, req.WithContext(ctx))

			statusKV := KV{K: HTTPStatusKey, V: rw.StatusCode}
			durKV := KV{K: HTTPDurationKey, V: timeSince(started).Milliseconds()}
			bytesKV := KV{K: HTTPBytesKey, V: rw.ContentLength}
			Info(ctx, KV{K: MessageKey, V: "end"}, methKV, urlKV, statusKV, durKV, bytesKV)
		})
	}
}

// Endpoint is a Goa endpoint middleware that adds the service and method names
// to the logged key/value pairs.
func Endpoint(e goa.Endpoint) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if s := ctx.Value(goa.ServiceKey); s != nil {
			ctx = With(ctx, KV{K: GoaServiceKey, V: s})
		}
		if m := ctx.Value(goa.MethodKey); m != nil {
			ctx = With(ctx, KV{K: GoaMethodKey, V: m})
		}
		return e(ctx, req)
	}
}

// Client wraps the given roundtripper and log requests and responses using the
// clue logger stored in the request context.
func Client(t http.RoundTripper, opts ...HTTPClientLogOption) http.RoundTripper {
	options := &httpClientOptions{
		iserr: func(status int) bool { return status >= 400 },
	}
	for _, o := range opts {
		o(options)
	}
	return &client{RoundTripper: t, options: options}
}

// WithPathFilter adds a path filter to the HTTP middleware. Requests whose path
// match the filter are not logged. WithPathFilter can be called multiple times
// to add multiple filters.
func WithPathFilter(filter *regexp.Regexp) HTTPLogOption {
	return func(o *httpLogOptions) {
		o.pathFilters = append(o.pathFilters, filter)
	}
}

// WithErrorStatus returns a HTTP client logger option that configures the
// logger to log errors for responses with the given status code.
func WithErrorStatus(status int) HTTPClientLogOption {
	return func(o *httpClientOptions) {
		o.iserr = func(s int) bool { return s == status }
	}
}

// WithLogBodyOnError returns a HTTP client logger option that configures the
// logger to log the response body when the response status code is an error.
func WithLogBodyOnError() HTTPClientLogOption {
	return func(o *httpClientOptions) {
		o.logErrBody = true
	}
}

// WithDisableRequestLogging returns a HTTP middleware option that disables
// logging of HTTP requests.
func WithDisableRequestLogging() HTTPLogOption {
	return func(o *httpLogOptions) {
		o.disableRequestLogging = true
	}
}

// WithDisableRequestID returns a HTTP middleware option that disables the
// generation of request IDs.
func WithDisableRequestID() HTTPLogOption {
	return func(o *httpLogOptions) {
		o.disableRequestID = true
	}
}

// RoundTrip executes the given HTTP request and logs the request and response. The
// request context must be initialized with a clue logger.
func (c *client) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	msgKV := KV{K: MessageKey, V: "finished client HTTP request"}
	methKV := KV{K: HTTPMethodKey, V: req.Method}
	urlKV := KV{K: HTTPURLKey, V: req.URL.String()}
	then := timeNow()
	resp, err = c.RoundTripper.RoundTrip(req)
	if err != nil {
		Error(req.Context(), err, msgKV, methKV, urlKV)
		return
	}
	ms := timeSince(then).Milliseconds()
	statusKV := KV{K: HTTPStatusKey, V: resp.Status}
	durKV := KV{K: HTTPDurationKey, V: ms}
	if c.options.iserr(resp.StatusCode) {
		if c.options.logErrBody {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				Error(req.Context(), err, msgKV, methKV, urlKV, statusKV, durKV)
				return resp, nil
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			Error(req.Context(), errors.New(resp.Status), msgKV, methKV, urlKV, statusKV, durKV, KV{K: HTTPBodyKey, V: string(body)})
		} else {
			Error(req.Context(), errors.New(resp.Status), msgKV, methKV, urlKV, statusKV, durKV)
		}
		return
	}
	Print(req.Context(), msgKV, methKV, urlKV, statusKV, durKV)
	return
}

// WriteHeader records the value of the status code before writing it.
func (w *responseCapture) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Write computes the written len and stores it in ContentLength.
func (w *responseCapture) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.ContentLength += n
	return n, err
}

// Flush implements the http.Flusher interface if the underlying response
// writer supports it.
func (w *responseCapture) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Push implements the http.Pusher interface if the underlying response
// writer supports it.
func (w *responseCapture) Push(target string, opts *http.PushOptions) error {
	if p, ok := w.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return errors.New("push not supported")
}

// Hijack supports the http.Hijacker interface.
func (w *responseCapture) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("response writer does not support hijacking: %T", w.ResponseWriter)
}

// from returns the client address from the request.
func from(req *http.Request) string {
	if f := req.Header.Get("X-Forwarded-For"); f != "" {
		return f
	}
	f := req.RemoteAddr
	ip, _, err := net.SplitHostPort(f)
	if err != nil {
		return f
	}
	return ip
}
