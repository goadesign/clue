package trace

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"goa.design/clue/log"
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
//      ctx := trace.Context(ctx, svcgen.ServiceName, trace.WithGRPCExporter(conn))
//      // Mount middleware
// 	handler := trace.HTTP(ctx)(mux)
//
func HTTP(ctx context.Context) func(http.Handler) http.Handler {
	s := ctx.Value(stateKey)
	if s == nil {
		panic(errContextMissing)
	}
	return func(h http.Handler) http.Handler {
		h = initTracingContext(ctx, h)
		h = addRequestIDHTTP(h)
		return otelhttp.NewHandler(h, s.(*stateBag).svc,
			otelhttp.WithTracerProvider(s.(*stateBag).provider),
			otelhttp.WithPropagators(s.(*stateBag).propagator))
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
		otelhttp.WithPropagators(s.(*stateBag).propagator))
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
func initTracingContext(traceCtx context.Context, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if IsTraced(ctx) {
			req = req.WithContext(withTracing(traceCtx, ctx))
			log.Debug(ctx,
				log.KV{log.TraceIDKey, trace.SpanFromContext(ctx).SpanContext().TraceID()})
		}
		h.ServeHTTP(w, req)
	})
}
