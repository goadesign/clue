package health

import (
	"context"
	"time"

	"goa.design/clue/log"
)

type (
	// Checker exposes a health check.
	Checker interface {
		// Check that all dependencies are healthy. Check returns true
		// if the service is healthy. The returned Health struct
		// contains the health status of each dependency.
		Check(context.Context) (*Health, bool)
	}

	// Health status of a service.
	Health struct {
		// Uptime of service in seconds.
		Uptime int64 `json:"uptime"`
		// Version of service.
		Version string `json:"version"`
		// Status of each dependency indexed by service name.
		// "OK" if dependency is healthy, "NOT OK" otherwise.
		Status map[string]string `json:"status,omitempty"`
	}

	// checker is a Checker that checks the health of the given
	// dependencies.
	checker struct {
		deps []Pinger
	}
)

// Version of service, initialized at compiled time.
var Version string

// StartedAt is the time the service was started.
var StartedAt = time.Now()

// Create a Checker that checks the health of the given dependencies.
func NewChecker(deps ...Pinger) Checker {
	return &checker{
		deps: deps,
	}
}

func (c *checker) Check(ctx context.Context) (*Health, bool) {
	res := &Health{
		Uptime:  int64(time.Since(StartedAt).Seconds()),
		Version: Version,
		Status:  make(map[string]string),
	}
	healthy := true
	for _, dep := range c.deps {
		res.Status[dep.Name()] = "OK"
		if err := dep.Ping(ctx); err != nil {
			res.Status[dep.Name()] = "NOT OK"
			healthy = false
			log.Error(ctx, err, log.KV{K: "msg", V: "ping failed"}, log.KV{K: "target", V: dep.Name()})
		}
	}
	return res, healthy
}
