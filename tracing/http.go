package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"goa.design/goa/v3/middleware"
)

// Middleware returns a tracing middleware that leverates the AWS Distro for
// OpenTelemetry to export traces to AWS X-Ray. It is aware of the Goa
// RequestID middleware and will use it to propagate the request ID to the
// trace.
//
// Example:
//
//	ctx := log.With(log.Context(context.Background()), "svc", svcgen.ServiceName)
//      svc := svc.New(ctx)
//      endpoints := svcgen.NewEndpoints(svc)
// 	mux := goa.NewMuxer()
//      httpsvr := httpsvrgen.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
//      httpsvrgen.Mount(mux, httpsvr)
//      provider := tracing.NewTracerProvider(ctx, svcgen.ServiceName, "localhost:6831")
// 	handler := tracing.Middleware(svcgen.ServiceName, provider)(mux)
// 	http.ListenAndServe(":8080", handler)
//
func Middleware(svc string, provider TracerProvider) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			h = addRequestIDHTTP(h)
			h = otelhttp.NewHandler(h, svc,
				otelhttp.WithTracerProvider(provider),
				otelhttp.WithPropagators(xray.Propagator{}),
				otelhttp.WithMeterProvider(metric.NewNoopMeterProvider()), // disable meter, use micro/instrument instead
			)
			h.ServeHTTP(w, req)
		})
	}
}

// addRequestIDHTTP is a middleware that adds the request ID to the current span
// attributes.
func addRequestIDHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := req.Context().Value(middleware.RequestIDKey)
		if requestID == nil {
			h.ServeHTTP(w, req)
			return
		}
		span := trace.SpanFromContext(req.Context())
		span.SetAttributes(attribute.String(AttributeRequestID, requestID.(string)))
		h.ServeHTTP(w, req)
	})
}
