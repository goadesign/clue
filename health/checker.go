package health

import (
	"context"
	"encoding/xml"
	"sort"
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

	// mp is used to marshal a map to xml.
	mp map[string]string
)

func (h Health) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(struct {
		XMLName xml.Name `xml:"health"`
		Uptime  int64    `xml:"uptime"`
		Version string   `xml:"version"`
		Status  mp       `xml:"status"`
	}{
		Uptime:  h.Uptime,
		Version: h.Version,
		Status:  h.Status,
	})
}

func (m mp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if err := e.EncodeElement(m[k], xml.StartElement{Name: xml.Name{Local: k}}); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

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
		// Note: need to create a new context for each dependency So that one
		// dependency canceling the context will not affect the other checks.
		logCtx := log.With(context.Background(), log.KV{K: "dep", V: dep.Name()})
		if err := dep.Ping(logCtx); err != nil {
			res.Status[dep.Name()] = "NOT OK"
			healthy = false
			log.Error(ctx, err, log.KV{K: "msg", V: "ping failed"})
		}
	}
	return res, healthy
}
