package metrics

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestOptions(t *testing.T) {
	var (
		durationBuckets     = []float64{1}
		requestSizeBuckets  = []float64{1}
		responseSizeBuckets = []float64{1}
		registerer          = NewTestRegistry(t)
		resolver            = func(_ *http.Request) string { return "test" }
	)
	options := defaultOptions()
	assertOptions(t, options, DefaultDurationBuckets, DefaultRequestSizeBuckets, DefaultResponseSizeBuckets, prometheus.DefaultRegisterer, nil)

	WithDurationBuckets(durationBuckets)(options)
	assertOptions(t, options, durationBuckets, DefaultRequestSizeBuckets, DefaultResponseSizeBuckets, prometheus.DefaultRegisterer, nil)

	WithRequestSizeBuckets(requestSizeBuckets)(options)
	assertOptions(t, options, durationBuckets, requestSizeBuckets, DefaultResponseSizeBuckets, prometheus.DefaultRegisterer, nil)

	WithResponseSizeBuckets(responseSizeBuckets)(options)
	assertOptions(t, options, durationBuckets, requestSizeBuckets, responseSizeBuckets, prometheus.DefaultRegisterer, nil)

	WithRegisterer(registerer)(options)
	assertOptions(t, options, durationBuckets, requestSizeBuckets, responseSizeBuckets, registerer, nil)

	WithRouteResolver(resolver)(options)
	assertOptions(t, options, durationBuckets, requestSizeBuckets, responseSizeBuckets, registerer, resolver)
}

func assertOptions(t *testing.T, options *options, durationBuckets []float64, requestSizeBuckets []float64, responseSizeBuckets []float64, registerer prometheus.Registerer, resolver RouteResolver) {
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
	if options.resolver == nil {
		if resolver != nil {
			t.Errorf("got nil, expected %v", resolver)
		}
		return
	}
	if resolver == nil {
		t.Errorf("got %v, expected nil", options.resolver)
		return
	}
	if fmt.Sprintf("%+v", options.resolver) != fmt.Sprintf("%+v", resolver) {
		t.Errorf("got %v, expected %v", options.resolver, resolver)
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
