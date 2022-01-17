package instrument

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"goa.design/goa/http/middleware"
	goa "goa.design/goa/v3/pkg"
)

// lengthReader is a wrapper around an io.ReadCloser that keeps track of how
// much data has been read.
type lengthReader struct {
	Source io.ReadCloser
	Length int
}

const (
	// MetricHTTPDuration is the name of the HTTP duration metric.
	MetricHTTPDuration = "http_server_duration_ms"
	// MetricHTTPRequestSize is the name of the HTTP request size metric.
	MetricHTTPRequestSize = "http_server_request_size_bytes"
	// MetricHTTPResponseSize is the name of the HTTP response size metric.
	MetricHTTPResponseSize = "http_server_response_size_bytes"
	// MetricHTTPActiveRequests is the name of the HTTP active requests metric_
	MetricHTTPActiveRequests = "http_server_active_requests"
	// LabelGoaService is the name of the label containing the Goa service name_
	LabelGoaService = "goa_service"
	// LabelHTTPVerb is the name of the label containing the HTTP verb_
	LabelHTTPVerb = "http_verb"
	// LabelHTTPHost is the name of the label containing the HTTP host_
	LabelHTTPHost = "http_host"
	// LabelHTTPPath is the name of the label containing the HTTP URL path_
	LabelHTTPPath = "http_path"
	// LabelHTTPStatusCode is the name of the label containing the HTTP status code_
	LabelHTTPStatusCode = "http_status_code"
)

var (
	// HTTPLabels is the set of dynamic labels used for all metrics but
	// MetricHTTPActiveRequests.
	HTTPLabels = []string{LabelHTTPVerb, LabelHTTPHost, LabelHTTPPath, LabelHTTPStatusCode}

	// HTTPActiveRequestsLabels is the set of dynamic labels used for the
	// MetricHTTPActiveRequests metric.
	HTTPActiveRequestsLabels = []string{LabelHTTPVerb, LabelHTTPHost, LabelHTTPPath}
)

// Be kind to tests
var timeSince = time.Since

// HTTP returns a middlware that collects the following metrics:
//
//    * `http.server.duration`: Histogram of request durations in milliseconds.
//    * `http.server.active_requests`: UpDownCounter of active requests.
//    * `http.server.request.size`: Histogram of request sizes in bytes.
//    * `http.server.response.size`: Histogram of response sizes in bytes.
//
// All the metrics have the following labels:
//
//    * `goa.method`: The method name as specified in the Goa design.
//    * `goa.service`: The service name as specified in the Goa design.
//    * `http.verb`: The HTTP verb (`GET`, `POST` etc.).
//    * `http.host`: The value of the HTTP host header.
//    * `http.path`: The HTTP path.
//    * `http.status_code`: The HTTP status code.
//
// Errors collecting or serving metrics are logged to the logger in the context
// if any.
func HTTP(svc string, opts ...Option) func(http.Handler) http.Handler {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}
	durations := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        MetricHTTPDuration,
		Help:        "Histogram of request durations in milliseconds.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
		Buckets:     options.durationBuckets,
	}, HTTPLabels)
	options.registerer.MustRegister(durations)

	reqSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        MetricHTTPRequestSize,
		Help:        "Histogram of request sizes in bytes.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
		Buckets:     options.requestSizeBuckets,
	}, HTTPLabels)
	options.registerer.MustRegister(reqSizes)

	respSizes := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        MetricHTTPResponseSize,
		Help:        "Histogram of response sizes in bytes.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
		Buckets:     options.responseSizeBuckets,
	}, HTTPLabels)
	options.registerer.MustRegister(respSizes)

	activeReqs := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        MetricHTTPActiveRequests,
		Help:        "Gauge of active requests.",
		ConstLabels: prometheus.Labels{LabelGoaService: svc},
	}, HTTPActiveRequestsLabels)
	options.registerer.MustRegister(activeReqs)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			labels := prometheus.Labels{
				LabelHTTPVerb: r.Method,
				LabelHTTPHost: r.Host,
				LabelHTTPPath: r.URL.Path,
			}
			activeReqs.With(labels).Add(1)
			defer activeReqs.With(labels).Sub(1)

			now := time.Now()
			rw := middleware.CaptureResponse(w)
			r.Body = &lengthReader{Source: r.Body}

			h.ServeHTTP(rw, r)

			labels[LabelHTTPStatusCode] = strconv.Itoa(rw.StatusCode)

			durations.With(labels).Observe(float64(timeSince(now).Milliseconds()))
			reqSizes.With(labels).Observe(float64(r.Body.(*lengthReader).Length))
			respSizes.With(labels).Observe(float64(rw.ContentLength))
		})
	}
}

func methodFromCtx(ctx context.Context) string {
	v := ctx.Value(goa.MethodKey)
	if v != nil {
		return v.(string)
	}
	return ""
}

func (r *lengthReader) Read(b []byte) (int, error) {
	n, err := r.Source.Read(b)
	r.Length += n
	return n, err
}

func (r *lengthReader) Close() error {
	var buf [32]byte
	var n int
	var err error
	for err == nil {
		n, err = r.Source.Read(buf[:])
		r.Length += n
	}
	closeerr := r.Source.Close()
	if err != nil && err != io.EOF {
		return err
	}
	return closeerr
}
