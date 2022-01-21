package tracing

import (
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestAdaptiveSampler(t *testing.T) {
	// We don't need to test Goa, keep it simple...
	s := adaptiveSampler(2, 10)
	expected := "Adaptive{maxSamplingRate:2,sampleSize:10}"
	if s.Description() != expected {
		t.Fatalf("got description %q, expected %q", s.Description(), expected)
	}
	res := s.ShouldSample(sdktrace.SamplingParameters{})
	if res.Decision != sdktrace.RecordAndSample {
		t.Error("expected sampling")
	}

	s2 := adaptiveSampler(1, 2)
	expected = "Adaptive{maxSamplingRate:1,sampleSize:2}"
	if s2.Description() != expected {
		t.Fatalf("got description %q, expected %q", s2.Description(), expected)
	}
	res2 := s2.ShouldSample(sdktrace.SamplingParameters{})
	res3 := s2.ShouldSample(sdktrace.SamplingParameters{})
	if res2.Decision != sdktrace.RecordAndSample {
		t.Error("expected sampling")
	}
	if res3.Decision != sdktrace.Drop {
		t.Error("expected no sampling")
	}
}
