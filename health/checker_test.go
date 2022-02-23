package health

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCheck(t *testing.T) {
	cases := []struct {
		name            string
		deps            []Pinger
		expectedStatus  map[string]string
		expectedHealthy bool
	}{
		{
			name:            "empty",
			expectedStatus:  map[string]string{},
			expectedHealthy: true,
		},
		{
			name:            "ok",
			deps:            singleHealthyDep("dependency"),
			expectedStatus:  map[string]string{"dependency": "OK"},
			expectedHealthy: true,
		},
		{
			name:            "not ok",
			deps:            singleUnhealthyDep("dependency", fmt.Errorf("dependency is not ok")),
			expectedStatus:  map[string]string{"dependency": "NOT OK"},
			expectedHealthy: false,
		},
		{
			name: "multiple dependencies",
			deps: multipleHealthyDeps("dependency1", "dependency2"),
			expectedStatus: map[string]string{
				"dependency1": "OK",
				"dependency2": "OK",
			},
			expectedHealthy: true,
		},
		{
			name: "multiple dependencies not ok",
			deps: multipleUnhealthyDeps(fmt.Errorf("dependency2 is not ok"), "dependency1", "dependency2"),
			expectedStatus: map[string]string{
				"dependency1": "OK",
				"dependency2": "NOT OK",
			},
			expectedHealthy: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			checker := NewChecker(c.deps...)
			res, err := checker.Check(context.Background())
			if err != c.expectedHealthy {
				t.Errorf("unexpected error: %v", err)
			}
			if res.Uptime != int64(time.Since(StartedAt).Seconds()) {
				t.Errorf("unexpected uptime: %d", res.Uptime)
			}
			if res.Version != Version {
				t.Errorf("unexpected version: %s", res.Version)
			}
			if len(res.Status) != len(c.expectedStatus) {
				t.Errorf("unexpected status: %v", res.Status)
			}
			for k, v := range c.expectedStatus {
				if res.Status[k] != v {
					t.Errorf("unexpected status for %s: %s", k, res.Status[k])
				}
			}
		})
	}
}

type mockDep struct {
	name string
	ping func(ctx context.Context) error
}

func (m *mockDep) Name() string                   { return m.name }
func (m *mockDep) Ping(ctx context.Context) error { return m.ping(ctx) }

func singleHealthyDep(name string) []Pinger {
	return []Pinger{&mockDep{
		name: name,
		ping: func(ctx context.Context) error {
			return nil
		},
	}}
}

func singleUnhealthyDep(name string, err error) []Pinger {
	return []Pinger{&mockDep{
		name: name,
		ping: func(ctx context.Context) error {
			return err
		},
	}}
}

func multipleHealthyDeps(names ...string) []Pinger {
	deps := make([]Pinger, len(names))
	for i, name := range names {
		deps[i] = &mockDep{
			name: name,
			ping: func(ctx context.Context) error {
				return nil
			},
		}
	}
	return deps
}

func multipleUnhealthyDeps(err error, names ...string) []Pinger {
	deps := make([]Pinger, len(names))
	for i, name := range names {
		deps[i] = &mockDep{
			name: name,
			ping: func(ctx context.Context) error {
				return nil
			},
		}
	}
	deps[len(deps)-1].(*mockDep).ping = func(ctx context.Context) error {
		return err
	}
	return deps
}
