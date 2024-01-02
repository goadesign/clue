package trace

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Message printed by panic when using a method with a non-initialized context.
const errContextMissing = "context not initialized for tracing, use trace.Context to set it up"

// HTTP returns a tracing middleware that uses an adaptive sampler for limiting
// the number of traced requests while guraranteeing that a certain number of
// requests are traced.
// HTTP panics if the context hasn't been initialized with Context.
//
// Example:
//
//	// Connect to remote trace collector.
//	conn, err := grpc.DialContext(ctx, collectorAddr)
//	if err != nil {
//	    log.Error(ctx, err)
//	    os.Exit(1)
//	}
//	// Initialize context for tracing
//	ctx := trace.Context(ctx, svcgen.ServiceName, trace.WithGRPCExporter(conn))
//	// Mount middleware
//	handler := trace.HTTP(ctx)(mux)
func HTTP(ctx context.Context) func(http.Handler) http.Handler {
	s := ctx.Value(stateKey)
	if s == nil {
		panic(errContextMissing)
	}
	return func(h http.Handler) http.Handler {
		return otelhttp.NewHandler(h, s.(*stateBag).svc,
			otelhttp.WithTracerProvider(s.(*stateBag).provider),
			otelhttp.WithPropagators(s.(*stateBag).propagator))
	}
}

// Client returns a http.RoundTripper that wraps t and creates spans for each
// request.  It panics if the context hasn't been initialized with Context.
//
// Example:
//
//	// Connect to remote trace collector.
//	conn, err := grpc.DialContext(ctx, collectorAddr)
//	if err != nil {
//	    log.Error(ctx, err)
//	    os.Exit(1)
//	}
//	// Initialize context for tracing
//	ctx := trace.Context(ctx, svcgen.ServiceName, trace.WithGRPCExporter(conn))
//	// Create client
//	cli := &http.Client{
//	    Transport: trace.Client(ctx, http.DefaultTransport),
//	}
//	// Use client
//	resp, err := cli.Get("http://example.com")
func Client(ctx context.Context, t http.RoundTripper, opts ...otelhttp.Option) http.RoundTripper {
	s := ctx.Value(stateKey)
	if s == nil {
		panic(errContextMissing)
	}
	opts = append(opts,
		otelhttp.WithTracerProvider(s.(*stateBag).provider),
		otelhttp.WithPropagators(s.(*stateBag).propagator))
	return otelhttp.NewTransport(t, opts...)
}
