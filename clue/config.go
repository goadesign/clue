package clue

import (
	"context"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"

	"goa.design/clue/log"

	// Force dependency on main module to ensure it is unambiguous during
	// module resolution.
	// See: https://github.com/googleapis/google-api-go-client/issues/2559.
	_ "google.golang.org/genproto/googleapis/type/datetime"
)

type (
	// Config is used to configure OpenTelemetry.
	Config struct {
		// MeterProvider is the OpenTelemetry meter provider used by clue
		MeterProvider metric.MeterProvider
		// TracerProvider is the OpenTelemetry tracer provider used clue
		TracerProvider trace.TracerProvider
		// Propagators is the OpenTelemetry propagator used by clue
		Propagators propagation.TextMapPropagator
		// ErrorHandler is the error handler used by OpenTelemetry
		ErrorHandler otel.ErrorHandler
	}
)

// ConfigureOpenTelemetry sets up code instrumentation using the OpenTelemetry
// API. It leverages the clue logger configured in ctx to log errors.
func ConfigureOpenTelemetry(ctx context.Context, cfg *Config) {
	otel.SetMeterProvider(cfg.MeterProvider)
	otel.SetTracerProvider(cfg.TracerProvider)
	otel.SetTextMapPropagator(cfg.Propagators)
	otel.SetLogger(logr.New(log.ToLogrSink(ctx)))
	otel.SetErrorHandler(cfg.ErrorHandler)
}

// NewConfig creates a new Config object adequate for use by
// ConfigureOpenTelemetry.  The metricExporter and spanExporter are used to
// record telemetry. If either is nil then the corresponding package will not
// record any telemetry. The OpenTelemetry metric provider is configured with a
// periodic reader. The OpenTelemetry tracer provider is configured to use a
// batch span processor and an adaptive sampler that aims at a maximum sampling
// rate of requests per second.  The resulting configuration can be modified
// (and providers replaced) by the caller prior to calling
// ConfigureOpenTelemetry.
//
// Example:
//
//	metricExporter, err := stdoutmetric.New()
//	if err != nil {
//		return err
//	}
//	spanExporter, err := stdouttrace.New()
//	if err != nil {
//		return err
//	}
//	cfg := clue.NewConfig(ctx, "mysvc", "1.0.0", metricExporter, spanExporter)
func NewConfig(
	ctx context.Context,
	svcName string,
	svcVersion string,
	metricExporter sdkmetric.Exporter,
	spanExporter sdktrace.SpanExporter,
	opts ...Option,
) (*Config, error) {
	options := defaultOptions(ctx)
	for _, o := range opts {
		o(options)
	}
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(svcName),
			semconv.ServiceVersionKey.String(svcVersion),
		))
	if err != nil {
		return nil, err
	}
	res, err = resource.Merge(res, options.resource)
	if err != nil {
		return nil, err
	}
	var meterProvider metric.MeterProvider
	if metricExporter == nil {
		meterProvider = metricnoop.NewMeterProvider()
	} else {
		var reader sdkmetric.Reader
		if options.readerInterval == 0 {
			reader = sdkmetric.NewPeriodicReader(metricExporter)
		} else {
			reader = sdkmetric.NewPeriodicReader(
				metricExporter,
				sdkmetric.WithInterval(options.readerInterval),
			)
		}
		meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(reader),
		)
	}
	var tracerProvider trace.TracerProvider
	if spanExporter == nil {
		tracerProvider = tracenoop.NewTracerProvider()
	} else {
		sampler := sdktrace.ParentBased(
			AdaptiveSampler(options.maxSamplingRate, options.sampleSize),
		)
		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
			sdktrace.WithSampler(sampler),
			sdktrace.WithBatcher(spanExporter),
		)
	}
	return &Config{
		MeterProvider:  meterProvider,
		TracerProvider: tracerProvider,
		Propagators:    options.propagators,
		ErrorHandler:   options.errorHandler,
	}, nil
}

// NewErrorHandler returns an error handler that logs errors using the clue
// logger configured in ctx.
func NewErrorHandler(ctx context.Context) otel.ErrorHandler {
	return otel.ErrorHandlerFunc(func(err error) {
		log.Error(ctx, err)
	})
}
