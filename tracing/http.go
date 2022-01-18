package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// HTTP returns a tracing middleware that leverates the AWS Distro for
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
// 	handler := tracing.HTTP(svcgen.ServiceName, provider)(mux)
// 	http.ListenAndServe(":8080", handler)
//
func HTTP(svc string, provider *sdktrace.TracerProvider) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return otelhttp.NewHandler(h, svc,
			otelhttp.WithTracerProvider(provider),
			otelhttp.WithPropagators(xray.Propagator{}),
		)
	}
}
