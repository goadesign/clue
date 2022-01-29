package metrics

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	svc := "testsvc"
	optionList := []Option{WithDurationBuckets([]float64{0, 1, 2})}
	optionStruct := defaultOptions()
	optionStruct.durationBuckets = []float64{0, 1, 2}
	otherOptionList := []Option{WithRequestSizeBuckets([]float64{0, 1, 2})}
	ctx := Context(context.Background(), svc, optionList...)
	cases := []struct {
		name    string
		ctx     context.Context
		svc     string
		options []Option
		want    stateBag
	}{
		{"empty", context.Background(), "", nil, stateBag{options: defaultOptions()}},
		{"with options", context.Background(), "", optionList, stateBag{options: optionStruct}},
		{"with options and service", context.Background(), svc, optionList, stateBag{options: optionStruct, svc: svc}},
		{"with initialized context", ctx, "some-svc", otherOptionList, stateBag{options: optionStruct, svc: svc}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := Context(c.ctx, c.svc, c.options...)
			got := ctx.Value(stateBagKey).(*stateBag)
			if got.svc != c.want.svc {
				t.Errorf("unexpected svc: got %q, want %q", got.svc, c.want.svc)
			}
			if !sameOptions(got.options, c.want.options) {
				t.Errorf("unexpected options: got %v, want %v", got.options, c.want.options)
			}
		})
	}
}

func sameOptions(a, b *options) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.durationBuckets != nil && b.durationBuckets != nil {
		if len(a.durationBuckets) != len(b.durationBuckets) {
			return false
		}
		for i := range a.durationBuckets {
			if a.durationBuckets[i] != b.durationBuckets[i] {
				return false
			}
		}
	}
	if a.requestSizeBuckets != nil && b.requestSizeBuckets != nil {
		if len(a.requestSizeBuckets) != len(b.requestSizeBuckets) {
			return false
		}
		for i := range a.requestSizeBuckets {
			if a.requestSizeBuckets[i] != b.requestSizeBuckets[i] {
				return false
			}
		}
	}
	return true
}
