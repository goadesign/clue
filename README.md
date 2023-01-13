# clue: Microservice Instrumentation

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![codecov](https://codecov.io/gh/goadesign/clue/branch/main/graph/badge.svg?token=HVP4WT1PS6)](https://codecov.io/gh/goadesign/clue)
[![Go Report Card](https://goreportcard.com/badge/goa.design/clue)](https://goreportcard.com/report/goa.design/clue) 
[![Go Reference](https://pkg.go.dev/badge/goa.design/clue.svg)](https://pkg.go.dev/goa.design/clue)

## Overview

`clue` provides a set of Go packages for instrumenting microservices. The
emphasis is on simplicity and ease of use. Although not a requirement, `clue`
works best when used in microservices written using
[Goa](https://github.com/goadesign/goa).

`clue` covers the following topics:

* Logging: the [log](log/) package provides a context-based logging API that
  intelligently selects what to log.
* Metrics: the [metrics](metrics/) package makes it possible for
  services to expose a Prometheus compatible `/metrics` HTTP endpoint.
* Health checks: the [health](health/) package provides a simple way for
  services to expose a health check endpoint.
* Dependency mocks: the [mock](mock/) package provides a way to mock
  downstream dependencies for testing.
* Tracing: the [trace](trace/) package conforms to the
  [OpenTelemetry](https://opentelemetry.io/) specification to trace requests.

The [weather](example/weather) example illustrates how to use `clue` to
instrument a system of Goa microservices. The example comes with a set of
scripts that can be used to compile and start the system as well as a complete
Grafana stack to query metrics and traces. See the
[README](example/weather/README.md) for more information.

## Contributing

See [Contributing](CONTRIBUTING.md)
