package tracing

import "testing"

func TestOptions(t *testing.T) {
	options := newDefaultOptions()
	if options.maxSamplingRate != 2 {
		t.Errorf("got %d, want 2", options.maxSamplingRate)
	}
	if options.sampleSize != 10 {
		t.Errorf("got %d, want 10", options.sampleSize)
	}
	WithMaxSamplingRate(3)(options)
	if options.maxSamplingRate != 3 {
		t.Errorf("got %d, want 3", options.maxSamplingRate)
	}
	WithSampleSize(20)(options)
	if options.sampleSize != 20 {
		t.Errorf("got %d, want 20", options.sampleSize)
	}
}
