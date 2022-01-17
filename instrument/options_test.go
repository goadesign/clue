package instrument

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestOptions(t *testing.T) {
	var (
		durationBuckets     = []float64{1}
		requestSizeBuckets  = []float64{1}
		responseSizeBuckets = []float64{1}
		registerer          = NewTestRegistry(t)
	)
	options := defaultOptions()
	assertOptions(t, options, DefaultDurationBuckets, DefaultRequestSizeBuckets, DefaultResponseSizeBuckets, prometheus.DefaultRegisterer)

	WithDurationBuckets(durationBuckets)(options)
	assertOptions(t, options, durationBuckets, DefaultRequestSizeBuckets, DefaultResponseSizeBuckets, prometheus.DefaultRegisterer)

	WithRequestSizeBuckets(requestSizeBuckets)(options)
	assertOptions(t, options, durationBuckets, requestSizeBuckets, DefaultResponseSizeBuckets, prometheus.DefaultRegisterer)

	WithResponseSizeBuckets(responseSizeBuckets)(options)
	assertOptions(t, options, durationBuckets, requestSizeBuckets, responseSizeBuckets, prometheus.DefaultRegisterer)

	WithRegisterer(registerer)(options)
	assertOptions(t, options, durationBuckets, requestSizeBuckets, responseSizeBuckets, registerer)
}

func assertOptions(t *testing.T, options *options, durationBuckets []float64, requestSizeBuckets []float64, responseSizeBuckets []float64, registerer prometheus.Registerer) {
	if !equal(options.durationBuckets, durationBuckets) {
		t.Errorf("got %v, expected %v", options.durationBuckets, durationBuckets)
	}
	if !equal(options.requestSizeBuckets, requestSizeBuckets) {
		t.Errorf("got %v, expected %v", options.requestSizeBuckets, requestSizeBuckets)
	}
	if !equal(options.responseSizeBuckets, responseSizeBuckets) {
		t.Errorf("got %v, expected %v", options.responseSizeBuckets, responseSizeBuckets)
	}
	if options.registerer != registerer {
		t.Errorf("got %v, expected %v", options.registerer, registerer)
	}
}

func equal(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
