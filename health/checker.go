package health

import (
	"context"
	"fmt"
	"time"

	"goa.design/clue/log"
)

type (
	// Checker exposes a health check.
	Checker interface {
		// Check that all dependencies are healthy. Check returns an
		// error if the service is unhealthy. The returned Health struct
		// contains the health status of each dependency.
		Check(context.Context) (*Health, error)
	}

	// Health status of a service.
	Health struct {
		// Uptime of service in seconds.
		Uptime int64 `json:"uptime"`
		// Git commit hash of service.
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

// Git commit hash of service, initialized at compiled time.
var GitCommit string

// StartedAt is the time the service was started.
var StartedAt = time.Now()

// Error returned when one ore more dependencies are unhealthy.
var ErrUnhealthy = fmt.Errorf("one or more dependencies are unhealthy")

// Create a Checker that checks the health of the given dependencies.
func NewChecker(deps ...Pinger) Checker {
	return &checker{
		deps: deps,
	}
}

func (c *checker) Check(ctx context.Context) (*Health, error) {
	res := &Health{
		Uptime:  int64(time.Since(StartedAt).Seconds()),
		Version: GitCommit,
		Status:  make(map[string]string),
	}
	healthy := true
	for _, dep := range c.deps {
		res.Status[dep.Name()] = "OK"
		if err := dep.Ping(ctx); err != nil {
			res.Status[dep.Name()] = "NOT OK"
			healthy = false
			log.Error(ctx, err, log.KV{"msg", "ping failed"}, log.KV{"target", dep.Name()})
		}
	}
	var err error
	if !healthy {
		err = ErrUnhealthy
	}
	return res, err
}
