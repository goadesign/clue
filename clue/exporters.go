package clue

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"

	"goa.design/clue/log"
)

// Allow mocking
var (
	otlpmetricgrpcNew = otlpmetricgrpc.New
	otlpmetrichttpNew = otlpmetrichttp.New
	otlptracegrpcNew  = otlptracegrpc.New
	otlptracehttpNew  = otlptracehttp.New
)

// NewGRPCMetricExporter returns an OpenTelementry Protocol metric exporter that
// report metrics to a gRPC collector.
func NewGRPCMetricExporter(ctx context.Context, options ...otlpmetricgrpc.Option) (exporter metric.Exporter, shutdown func(), err error) {
	exporter, err = otlpmetricgrpcNew(ctx, options...)
	if err != nil {
		return
	}
	shutdown = func() {
		// Create new context in case the parent context has been canceled.
		ctx := log.WithContext(context.Background(), ctx)
		if err := exporter.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown metric exporter")
		}
	}
	return
}

// NewGRPCSpanExporter returns an OpenTelementry Protocol span exporter that
// report spans to a gRPC collector.
func NewGRPCSpanExporter(ctx context.Context, options ...otlptracegrpc.Option) (exporter trace.SpanExporter, shutdown func(), err error) {
	exporter, err = otlptracegrpcNew(ctx, options...)
	if err != nil {
		return
	}
	shutdown = func() {
		// Create new context in case the parent context has been canceled.
		ctx := log.WithContext(context.Background(), ctx)
		if err := exporter.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown span exporter")
		}
	}
	return
}

// NewHTTPMetricExporter returns an OpenTelementry Protocol metric exporter that
// report metrics to a HTTP collector.
func NewHTTPMetricExporter(ctx context.Context, options ...otlpmetrichttp.Option) (exporter metric.Exporter, shutdown func(), err error) {
	exporter, err = otlpmetrichttpNew(ctx, options...)
	if err != nil {
		return
	}
	shutdown = func() {
		// Create new context in case the parent context has been canceled.
		ctx := log.WithContext(context.Background(), ctx)
		if err := exporter.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown metric exporter")
		}
	}
	return
}

// NewHTTPSpanExporter returns an OpenTelementry Protocol span exporter that
// report spans to a HTTP collector.
func NewHTTPSpanExporter(ctx context.Context, options ...otlptracehttp.Option) (exporter trace.SpanExporter, shutdown func(), err error) {
	exporter, err = otlptracehttpNew(ctx, options...)
	if err != nil {
		return
	}
	shutdown = func() {
		// Create new context in case the parent context has been canceled.
		ctx := log.WithContext(context.Background(), ctx)
		if err := exporter.Shutdown(ctx); err != nil {
			log.Errorf(ctx, err, "failed to shutdown span exporter")
		}
	}
	return
}
