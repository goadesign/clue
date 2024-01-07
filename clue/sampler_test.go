package clue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestAdaptiveSampler(t *testing.T) {
	// We don't need to test Goa, keep it simple...
	s := AdaptiveSampler(2, 10)
	expected := "Adaptive{maxSamplingRate:2,sampleSize:10}"
	assert.Equal(t, expected, s.Description())
	res := s.ShouldSample(sdktrace.SamplingParameters{})
	assert.Equal(t, sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample}, res)
	s2 := AdaptiveSampler(1, 2)
	expected = "Adaptive{maxSamplingRate:1,sampleSize:2}"
	assert.Equal(t, expected, s2.Description())
	res2 := s2.ShouldSample(sdktrace.SamplingParameters{})
	res3 := s2.ShouldSample(sdktrace.SamplingParameters{})
	assert.Equal(t, sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample}, res2)
	assert.Equal(t, sdktrace.SamplingResult{Decision: sdktrace.Drop}, res3)
}
