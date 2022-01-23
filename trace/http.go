package trace

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"goa.design/goa/v3/middleware"
)

// Message printed by panic when using a method with a non-initialized context.
const errContextMissing = "context not initialized for tracing, use trace.Context to set it up"

// HTTP returns a tracing middleware that uses a parent based sampler (i.e.
// traces if the parent request traces) and an adaptive root sampler (i.e.  when
// there is no parent uses a target number of requests per second to trace).
// The implementation leverages the OpenTelemetry SDK and can thus be configured
// to send traces to an OpenTelemetry remote collector. It is aware of the Goa
// RequestID middleware and will use it to propagate the request ID to the
// trace. HTTP panics if the context hasn't been initialized with Context.
//
// Example:
//
//      // Connect to remote trace collector.
//      conn, err := grpc.DialContext(ctx, collectorAddr,
//          grpc.WithTransportCrendentials(insecure.Credentials()))
//      if err != nil {
//          log.Error(ctx, err)
//          os.Exit(1)
//      }
//      // Initialize context for tracing
//      ctx := trace.Context(ctx, svcgen.ServiceName, conn)
//      // Mount middleware
// 	handler := trace.HTTP(ctx, svcgen.ServiceName)(mux)
//
func HTTP(ctx context.Context, svc string) func(http.Handler) http.Handler {
	s := ctx.Value(stateKey)
	if s == nil {
		panic(errContextMissing)
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			h = initTracingContext(ctx, h)
			h = addRequestIDHTTP(h)
			h = otelhttp.NewHandler(h, svc,
				otelhttp.WithTracerProvider(s.(*stateBag).provider),
				// disable meter, use micro/instrument instead
				otelhttp.WithMeterProvider(metric.NewNoopMeterProvider()),
			)
			h.ServeHTTP(w, req)
		})
	}
}

// Client returns a roundtripper that wraps t and creates spans for each request.
// It panics if the context hasn't been initialized with Context.
func Client(ctx context.Context, t http.RoundTripper) http.RoundTripper {
	s := ctx.Value(stateKey)
	if s == nil {
		panic(errContextMissing)
	}
	return otelhttp.NewTransport(t,
		otelhttp.WithTracerProvider(s.(*stateBag).provider),
		otelhttp.WithMeterProvider(metric.NewNoopMeterProvider()))
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

// initTracingContext is a middleware that adds the tracing state to the request
// context.
func initTracingContext(ctx context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s := ctx.Value(stateKey).(*stateBag)
		ctx := withProvider(req.Context(), s.provider)
		setActiveSpans(ctx, []trace.Span{trace.SpanFromContext(req.Context())})
		h.ServeHTTP(w, req.WithContext(ctx))
	})
}