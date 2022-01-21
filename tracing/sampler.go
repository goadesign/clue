package tracing

import (
	"fmt"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"goa.design/goa/v3/middleware"
)

// sampler leverages the Goa adaptive sampler implementation.
type sampler struct {
	s               middleware.Sampler
	maxSamplingRate int
	sampleSize      int
}

// adaptiveSampler computes the interval for sampling for tracing middleware.
// it can also be used by non-web go routines to trace internal API calls.
//
// maxSamplingRate is the desired maximum sampling rate in requests per second.
//
// sampleSize sets the number of requests between two adjustments of the
// sampling rate when MaxSamplingRate is set. the sample rate cannot be adjusted
// until the sample size is reached at least once.
func adaptiveSampler(maxSamplingRate, sampleSize int) sdktrace.Sampler {
	return sampler{
		s:               middleware.NewAdaptiveSampler(maxSamplingRate, sampleSize),
		maxSamplingRate: maxSamplingRate,
		sampleSize:      sampleSize,
	}
}

// Description returns the description of the sampler.
func (s sampler) Description() string {
	return fmt.Sprintf("Adaptive{maxSamplingRate:%d,sampleSize:%d}", s.maxSamplingRate, s.sampleSize)
}

// ShouldSample returns the sampling decision for the given parameters.
func (s sampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if !s.s.Sample() {
		return sdktrace.SamplingResult{Decision: sdktrace.Drop}
	}

	psc := trace.SpanContextFromContext(p.ParentContext)
	return sdktrace.SamplingResult{
		Decision:   sdktrace.RecordAndSample,
		Tracestate: psc.TraceState(),
	}
}
