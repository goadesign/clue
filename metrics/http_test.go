package metrics

import (
	"context"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"goa.design/clue/internal/testsvc"
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg), WithDurationBuckets(buckets))
			middleware := HTTP(ctx, nil)
			cli, stop := testsvc.SetupHTTP(t,
				testsvc.WithHTTPMiddleware(middleware),
				testsvc.WithHTTPFunc(noopMethod()))
			_, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(metricHTTPDuration, httpLabels, 1, c.expectedBucketCounts)
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg), WithRequestSizeBuckets(buckets))
			middleware := HTTP(ctx, nil)
			cli, stop := testsvc.SetupHTTP(t,
				testsvc.WithHTTPMiddleware(middleware),
				testsvc.WithHTTPFunc(noopMethod()))

			_, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{S: &c.str})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(metricHTTPRequestSize, httpLabels, 1, c.expectedBucketCounts)
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
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg), WithResponseSizeBuckets(buckets))
			middleware := HTTP(ctx, nil)
			cli, stop := testsvc.SetupHTTP(t,
				testsvc.WithHTTPMiddleware(middleware),
				testsvc.WithHTTPFunc(stringMethod(c.str)))

			_, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			stop()
			reg.AssertHistogram(metricHTTPResponseSize, httpLabels, 1, c.expectedBucketCounts)
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
		{"one hundred", 100},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reg := NewTestRegistry(t)
			ctx := Context(context.Background(), "testsvc", WithRegisterer(reg))
			middleware := HTTP(ctx, nil)
			chstop := make(chan struct{})
			var running, done sync.WaitGroup
			running.Add(c.numReqs)
			done.Add(c.numReqs)
			cli, stop := testsvc.SetupHTTP(t,
				testsvc.WithHTTPMiddleware(middleware),
				testsvc.WithHTTPFunc(waitMethod(&running, &done, chstop)))

			for i := 0; i < c.numReqs; i++ {
				go func() {
					_, err := cli.HTTPMethod(context.Background(), &testsvc.Fields{})
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
				}()
			}

			running.Wait()
			reg.AssertGauge(metricHTTPActiveRequests, httpActiveRequestsLabels, c.numReqs)
			close(chstop)
			done.Wait()
			stop()
		})
	}
}

func TestLengthReader(t *testing.T) {
	cases := []struct {
		name         string
		str          string
		expectedSize int
	}{
		{"empty", "", 0},
		{"one", "1", 1},
		{"ten", strings.Repeat("1", 10), 10},
		{"one hundred", strings.Repeat("1", 100), 100},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := strings.NewReader(c.str)
			ctx, lr := newLengthReader(io.NopCloser(r), context.Background())
			n, err := lr.Read(make([]byte, 100))
			if err != nil && err != io.EOF {
				t.Errorf("unexpected error: %v", err)
			}
			if n != c.expectedSize {
				t.Errorf("expected %d bytes, got %d", c.expectedSize, n)
			}
			length := ctx.Value(ctxReqLen)
			if length == nil {
				t.Fatal("expected length to be set in context")
			}
			if *(length.(*int)) != c.expectedSize {
				t.Errorf("expected %d bytes, got %d", c.expectedSize, *(length.(*int)))
			}
			err = lr.Close()
			if err != nil {
				t.Errorf("unexpected close error: %v", err)
			}
		})
	}
}
