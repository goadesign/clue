package clue

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"goa.design/clue/log"
)

func TestNewGRPCMetricExporter(t *testing.T) {
	testErr := errors.New("test error")
	// Define test cases
	tests := []struct {
		name    string
		options []otlpmetricgrpc.Option
		newErr  error
		wantLog string
		wantErr bool
	}{
		{
			name: "Success",
		},
		{
			name:    "Options",
			options: []otlpmetricgrpc.Option{otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint("test")},
		},
		{
			name:    "New Error",
			newErr:  testErr,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otlpmetricgrpcNew = func(ctx context.Context, _ ...otlpmetricgrpc.Option) (*otlpmetricgrpc.Exporter, error) {
				if tt.newErr != nil {
					return nil, tt.newErr
				}
				return otlpmetricgrpc.New(ctx)
			}
			var buf bytes.Buffer
			ctx := log.Context(context.Background(), log.WithOutputs(log.Output{Writer: &buf, Format: log.FormatText}))

			exporter, shutdown, err := NewGRPCMetricExporter(ctx, tt.options...)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, exporter)
				assert.Nil(t, shutdown)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, exporter)
			assert.NotNil(t, shutdown)
			shutdown()
			assert.Empty(t, buf.String())
		})
	}
}

func TestNewHTTPMetricExporter(t *testing.T) {
	testErr := errors.New("test error")
	// Define test cases
	tests := []struct {
		name    string
		options []otlpmetrichttp.Option
		newErr  error
		wantLog string
		wantErr bool
	}{
		{
			name: "Success",
		},
		{
			name:    "Options",
			options: []otlpmetrichttp.Option{otlpmetrichttp.WithEndpoint("test")},
		},
		{
			name:    "New Error",
			newErr:  testErr,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otlpmetrichttpNew = func(ctx context.Context, _ ...otlpmetrichttp.Option) (*otlpmetrichttp.Exporter, error) {
				if tt.newErr != nil {
					return nil, tt.newErr
				}
				return otlpmetrichttp.New(ctx)
			}
			var buf bytes.Buffer
			ctx := log.Context(context.Background(), log.WithOutputs(log.Output{Writer: &buf, Format: log.FormatText}))

			exporter, shutdown, err := NewHTTPMetricExporter(ctx, tt.options...)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, exporter)
				assert.Nil(t, shutdown)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, exporter)
			assert.NotNil(t, shutdown)
			shutdown()
			assert.Empty(t, buf.String())
		})
	}
}

func TestNewGRPCSpanExporter(t *testing.T) {
	testErr := errors.New("test error")
	// Define test cases
	tests := []struct {
		name    string
		options []otlptracegrpc.Option
		newErr  error
		wantLog string
		wantErr bool
	}{
		{
			name: "Success",
		},
		{
			name:    "Options",
			options: []otlptracegrpc.Option{otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("test")},
		},
		{
			name:    "New Error",
			newErr:  testErr,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otlptracegrpcNew = func(ctx context.Context, _ ...otlptracegrpc.Option) (*otlptrace.Exporter, error) {
				if tt.newErr != nil {
					return nil, tt.newErr
				}
				return otlptracegrpc.New(ctx)
			}
			var buf bytes.Buffer
			ctx := log.Context(context.Background(), log.WithOutputs(log.Output{Writer: &buf, Format: log.FormatText}))

			exporter, shutdown, err := NewGRPCSpanExporter(ctx, tt.options...)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, exporter)
				assert.Nil(t, shutdown)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, exporter)
			assert.NotNil(t, shutdown)
			shutdown()
			assert.Empty(t, buf.String())
		})
	}
}

func TestNewHTTPSpanExporter(t *testing.T) {
	testErr := errors.New("test error")
	// Define test cases
	tests := []struct {
		name    string
		options []otlptracehttp.Option
		newErr  error
		wantLog string
		wantErr bool
	}{
		{
			name: "Success",
		},
		{
			name:    "Options",
			options: []otlptracehttp.Option{otlptracehttp.WithEndpoint("test")},
		},
		{
			name:    "New Error",
			newErr:  testErr,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otlptracehttpNew = func(ctx context.Context, _ ...otlptracehttp.Option) (*otlptrace.Exporter, error) {
				if tt.newErr != nil {
					return nil, tt.newErr
				}
				return otlptracehttp.New(ctx)
			}
			var buf bytes.Buffer
			ctx := log.Context(context.Background(), log.WithOutputs(log.Output{Writer: &buf, Format: log.FormatText}))

			exporter, shutdown, err := NewHTTPSpanExporter(ctx, tt.options...)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, exporter)
				assert.Nil(t, shutdown)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, exporter)
			assert.NotNil(t, shutdown)
			shutdown()
			assert.Empty(t, buf.String())
		})
	}
}
