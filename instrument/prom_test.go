package instrument

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

type Registry struct {
	*prometheus.Registry
	t *testing.T
}

var _ prometheus.Registerer = (*Registry)(nil)

func NewTestRegistry(t *testing.T) *Registry {
	return &Registry{prometheus.NewRegistry(), t}
}

// AssertGauge validates that the gauge with the given name and labels exists
// and has the given value.
func (r *Registry) AssertGauge(name string, labels []string, value int) {
	metric := r.findMetric(name, labels)
	if metric.Gauge == nil {
		r.t.Errorf("gauge %q with labels %v not found", name, labels)
		return
	}
	var val float64
	if metric.Gauge.Value != nil {
		val = *metric.Gauge.Value
	}
	if float64(value) != val {
		r.t.Errorf("gauge %q with labels %v has value %v, want %v", name, labels, val, value)
	}
}

// AssertHistogram validates that the histogram with the given name and given
// labels exists and has the given sample count and cumulative counts for the
// given buckets.
func (r *Registry) AssertHistogram(name string, labels []string, sampleCount int, bucketCumulativeCount []int) {
	metric := r.findMetric(name, labels)
	if metric.Histogram == nil {
		r.t.Errorf("histogram %s with labels %v not found", name, labels)
		return
	}
	count := 0
	if metric.Histogram.SampleCount != nil {
		count = int(*metric.Histogram.SampleCount)
	}
	if count != sampleCount {
		r.t.Errorf("histogram %s with labels %v has sample count %d, want %d", name, labels, count, sampleCount)
	}
	if len(metric.Histogram.Bucket) != len(bucketCumulativeCount) {
		r.t.Fatalf("histogram %s with labels %v has %d buckets, want %d", name, labels, len(metric.Histogram.Bucket), len(bucketCumulativeCount))
	}
	for i, b := range metric.Histogram.Bucket {
		count := 0
		if b.CumulativeCount != nil {
			count = int(*b.CumulativeCount)
		}
		if count != bucketCumulativeCount[i] {
			r.t.Errorf("histogram %s with labels %v has bucket %d cumulative count %d, want %d", name, labels, i, count, bucketCumulativeCount[i])
		}
	}
}

// findMetric finds a metric in the registry with the given name and labels.
func (r *Registry) findMetric(name string, labels []string) *io_prometheus_client.Metric {
	families, err := r.Gather()
	if err != nil {
		r.t.Errorf("failed to gather metrics: %v", err)
	}
	var metrics *io_prometheus_client.Metric
loop:
	for _, family := range families {
		if family.Name == nil || *family.Name != name {
			continue
		}
		for _, m := range family.Metric {
			if !hasLabels(m.Label, labels) {
				continue
			}
			metrics = m
			break loop
		}
	}
	if metrics == nil {
		r.t.Fatalf("histogram %s with labels %v not found", name, labels)
	}
	return metrics
}

// hasLabels returns true if the given label names are a subset of the given
// prometheus label names.
func hasLabels(promlabels []*io_prometheus_client.LabelPair, names []string) bool {
	for _, lbl := range names {
		found := false
		for _, plbl := range promlabels {
			if plbl.Name == nil {
				continue
			}
			if *plbl.Name == lbl {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
