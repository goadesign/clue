package metrics

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"goa.design/goa/v3/http/middleware"
)

type (
	// lengthReader is a wrapper around an io.ReadCloser that keeps track of how
	// much data has been read.
	lengthReader struct {
		Source io.ReadCloser
		ctx    context.Context
	}
)

// Be kind to tests
var timeSince = time.Since

// HTTP returns a middlware that metricss requests. The context must have
// been initialized with Context. HTTP collects the following metrics:
//
//   - `http.server.duration`: Histogram of request durations in milliseconds.
//   - `http.server.active_requests`: UpDownCounter of active requests.
//   - `http.server.request.size`: Histogram of request sizes in bytes.
//   - `http.server.response.size`: Histogram of response sizes in bytes.
//
// All the metrics have the following labels:
//
//   - `http.verb`: The HTTP verb (`GET`, `POST` etc.).
//   - `http.host`: The value of the HTTP host header.
//   - `http.path`: The HTTP path.
//   - `http.status_code`: The HTTP status code.
//
// Errors collecting or serving metrics are logged to the logger in the context
// if any.
func HTTP(ctx context.Context) func(http.Handler) http.Handler {
	b := ctx.Value(stateBagKey)
	if b == nil {
		panic("initialize context with Context first")
	}
	metrics := b.(*stateBag).HTTPMetrics()
	resolver := b.(*stateBag).options.resolver

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var route string
			if resolver != nil {
				route = resolver(req)
			} else {
				route = req.URL.Path
			}
			labels := prometheus.Labels{
				labelHTTPVerb: req.Method,
				labelHTTPHost: req.Host,
				labelHTTPPath: route,
			}
			metrics.ActiveRequests.With(labels).Add(1)
			defer metrics.ActiveRequests.With(labels).Sub(1)

			now := time.Now()
			rw := middleware.CaptureResponse(w)
			ctx, body := newLengthReader(req.Body, req.Context())
			req.Body = body
			req = req.WithContext(ctx)

			h.ServeHTTP(rw, req)

			labels[labelHTTPStatusCode] = strconv.Itoa(rw.StatusCode)

			reqLength := req.Context().Value(ctxReqLen).(*int)
			metrics.Durations.With(labels).Observe(float64(timeSince(now).Milliseconds()))
			metrics.RequestSizes.With(labels).Observe(float64(*reqLength))
			metrics.ResponseSizes.With(labels).Observe(float64(rw.ContentLength))
		})
	}
}

// So we have to do a little dance to get the length of the request body.  We
// can't just simply wrap the body and sum up the length on each read because
// otel sets its own wrapper which means we can't cast the request back after
// the call to the next handler. We thus store the computed length in the
// context instead.
func newLengthReader(body io.ReadCloser, ctx context.Context) (context.Context, *lengthReader) {
	reqLen := 0
	ctx = context.WithValue(ctx, ctxReqLen, &reqLen)
	return ctx, &lengthReader{body, ctx}
}

func (r *lengthReader) Read(b []byte) (int, error) {
	n, err := r.Source.Read(b)
	l := r.ctx.Value(ctxReqLen).(*int)
	*l += n

	return n, err
}

func (r *lengthReader) Close() error {
	var buf [32]byte
	var n int
	var err error
	for err == nil {
		n, err = r.Source.Read(buf[:])
		l := r.ctx.Value(ctxReqLen).(*int)
		*l += n
	}
	closeerr := r.Source.Close()
	if err != nil && err != io.EOF {
		return err
	}
	return closeerr
}
