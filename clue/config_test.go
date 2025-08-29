package clue

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	lognoop "go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	"goa.design/clue/log"
)

type (
	// dummyErrorHandler is a dummy implementation of the OpenTelemetry error handler interface.
	dummyErrorHandler struct{}
)

func TestConfigureOpenTelemetry(t *testing.T) {
	ctx := log.Context(context.Background())
	noopMeterProvider := metricnoop.NewMeterProvider()
	noopTracerProvider := tracenoop.NewTracerProvider()
	noopLoggerProvider := lognoop.NewLoggerProvider()
	noopErrorHandler := dummyErrorHandler{}

	cases := []struct {
		name           string
		meterProvider  metric.MeterProvider
		tracerProvider trace.TracerProvider
		loggerProvider otellog.LoggerProvider
		propagators    propagation.TextMapPropagator
		errorHandler   otel.ErrorHandler

		wantMeterProvider  metric.MeterProvider
		wantTracerProvider trace.TracerProvider
		wantLoggerProvider otellog.LoggerProvider
		wantPropagators    propagation.TextMapPropagator
		wantErrorHandler   bool
	}{
		{
			name: "default",
		}, {
			name:              "meter provider",
			meterProvider:     noopMeterProvider,
			wantMeterProvider: noopMeterProvider,
		}, {
			name:               "tracer provider",
			tracerProvider:     noopTracerProvider,
			wantTracerProvider: noopTracerProvider,
		}, {
			name:               "logger provider",
			loggerProvider:     noopLoggerProvider,
			wantLoggerProvider: noopLoggerProvider,
		}, {
			name:            "propagators",
			propagators:     propagation.Baggage{},
			wantPropagators: propagation.Baggage{},
		}, {
			name:             "error handler",
			errorHandler:     &noopErrorHandler,
			wantErrorHandler: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := &Config{
				MeterProvider:  c.meterProvider,
				TracerProvider: c.tracerProvider,
				LoggerProvider: c.loggerProvider,
				Propagators:    c.propagators,
				ErrorHandler:   c.errorHandler,
			}
			ConfigureOpenTelemetry(ctx, cfg)
			assert.Equal(t, c.wantMeterProvider, otel.GetMeterProvider())
			assert.Equal(t, c.wantTracerProvider, otel.GetTracerProvider())
			assert.Equal(t, c.wantLoggerProvider, global.GetLoggerProvider())
			assert.Equal(t, c.wantPropagators, otel.GetTextMapPropagator())
		})
	}
}

func TestNewConfig(t *testing.T) {
	ctx := log.Context(context.Background())
	svcName := "svcName"
	svcVersion := "svcVersion"
	spanExporter, err := stdouttrace.New()
	require.NoError(t, err)
	metricsExporter, err := stdoutmetric.New()
	require.NoError(t, err)
	logExporter, err := stdoutlog.New()
	require.NoError(t, err)
	noopErrorHandler := dummyErrorHandler{}

	cases := []struct {
		name            string
		metricsExporter sdkmetric.Exporter
		spanExporter    sdktrace.SpanExporter
		logExporter     sdklog.Exporter
		propagators     propagation.TextMapPropagator
		errorHandler    otel.ErrorHandler

		wantPropagators  propagation.TextMapPropagator
		wantErrorHandler bool
	}{
		{
			name: "default",
		}, {
			name:            "metrics exporter",
			metricsExporter: metricsExporter,
		}, {
			name:         "tracer provider",
			spanExporter: spanExporter,
		}, {
			name:        "log exporter",
			logExporter: logExporter,
		}, {
			name:            "propagators",
			propagators:     propagation.Baggage{},
			wantPropagators: propagation.Baggage{},
		}, {
			name:             "error handler",
			errorHandler:     &noopErrorHandler,
			wantErrorHandler: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg, err := NewConfig(ctx,
				svcName,
				svcVersion,
				c.metricsExporter,
				c.spanExporter,
				c.logExporter,
				WithPropagators(c.propagators),
				WithErrorHandler(c.errorHandler))
			require.NoError(t, err)
			if c.spanExporter != nil {
				serialized := fmt.Sprintf("%+v", cfg.TracerProvider)
				assert.Contains(t, serialized, "maxSamplingRate:2")
			}
			assert.Equal(t, c.wantPropagators, cfg.Propagators)
			assert.Equal(t, c.wantErrorHandler, cfg.ErrorHandler != nil)
		})
	}
}

func (dummyErrorHandler) Handle(error) {}
