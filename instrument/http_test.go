package instrument

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/crossnokaye/micro/instrument/testsvc"
)

func TestHTTPServerDuration(t *testing.T) {
	buckets := []float64{10, 110}
	cases := []struct {
		name                 string
		d                    time.Duration
		expectedBucketCounts []int
	}{
		{"fast", 1 * time.Millisecond, []int{1, 1}},
		{"slow", 100 * time.Millisecond, []int{0, 1}},
		{"very slow", 1000 * time.Millisecond, []int{0, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			restore := timeSince
			defer func() { timeSince = restore }()
			timeSince = func(time.Time) time.Duration { return c.d }

			reg := NewTestRegistry(t)
			middleware := HTTP("testsvc", WithRegisterer(reg), WithDurationBuckets(buckets))
			cli, stop := testsvc.SetupHTTP(t,
				testsvc.WithHTTPMiddleware(middleware),
				testsvc.WithHTTPFunc(noopUnaryMethod()))

			_, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(MetricHTTPDuration, HTTPLabels, 1, c.expectedBucketCounts)
		})
	}
}

func TestHTTPRequestSize(t *testing.T) {
	buckets := []float64{10, 110}
	cases := []struct {
		name                 string
		str                  string
		expectedBucketCounts []int
	}{
		{"small", "1", []int{1, 1}},
		{"large", strings.Repeat("1", 100), []int{0, 1}},
		{"very large", strings.Repeat("1", 1000), []int{0, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := NewTestRegistry(t)
			middleware := HTTP("testsvc", WithRegisterer(reg), WithRequestSizeBuckets(buckets))
			cli, stop := testsvc.SetupHTTP(t,
				testsvc.WithHTTPMiddleware(middleware),
				testsvc.WithHTTPFunc(noopUnaryMethod()))

			_, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{S: &c.str})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(MetricHTTPRequestSize, HTTPLabels, 1, c.expectedBucketCounts)
		})
	}
}

func TestHTTPResponseSize(t *testing.T) {
	buckets := []float64{10, 110}
	cases := []struct {
		name                 string
		str                  string
		expectedBucketCounts []int
	}{
		{"small", "1", []int{1, 1}},
		{"large", strings.Repeat("1", 100), []int{0, 1}},
		{"very large", strings.Repeat("1", 1000), []int{0, 0}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := NewTestRegistry(t)
			middleware := HTTP("testsvc", WithRegisterer(reg), WithResponseSizeBuckets(buckets))
			cli, stop := testsvc.SetupHTTP(t,
				testsvc.WithHTTPMiddleware(middleware),
				testsvc.WithHTTPFunc(stringUnaryMethod(c.str)))

			_, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(MetricHTTPResponseSize, HTTPLabels, 1, c.expectedBucketCounts)
		})
	}
}

func TestHTTPActiveRequests(t *testing.T) {
	cases := []struct {
		name    string
		numReqs int
	}{
		{"one", 1},
		{"ten", 10},
		{"one thousand", 1000},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := NewTestRegistry(t)
			middleware := HTTP("testsvc", WithRegisterer(reg))
			chstop := make(chan struct{})
			var running, done sync.WaitGroup
			running.Add(c.numReqs)
			done.Add(c.numReqs)
			cli, stop := testsvc.SetupHTTP(t,
				testsvc.WithHTTPMiddleware(middleware),
				testsvc.WithHTTPFunc(waitUnaryMethod(&running, &done, chstop)))

			for i := 0; i < c.numReqs; i++ {
				go func() {
					_, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{})
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
				}()
			}

			running.Wait()
			reg.AssertGauge(MetricHTTPActiveRequests, HTTPActiveRequestsLabels, c.numReqs)
			close(chstop)
			done.Wait()
			reg.AssertGauge(MetricHTTPActiveRequests, HTTPActiveRequestsLabels, 0)
			stop()
		})
	}
}
