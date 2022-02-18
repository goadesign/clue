# clue: Microservice Instrumentation

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![codecov](https://codecov.io/gh/goadesign/clue/branch/main/graph/badge.svg?token=HVP4WT1PS6)](https://codecov.io/gh/goadesign/clue)

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

## Roadmap to Open Source

The goal is to make this repo public ASAP. The one remaining item is a rework of
the `log` package to make it more flexible and allow both `Printf` and key-based
style logging.

## Addendum: Importing Private Repository

Because it is currently private importing this repo requires a few extra steps:

1. The `go get` command must be able to use `ssh` to clone repoitories hosted on
   GitHub.
2. The `go get` command must know that this repo is private so it uses the right
   `git` command.
   
These two requirements can be satisfied with the corresponding two commands below:

```bash
git config --global url.git@github.com:.insteadOf https://github.com/
go env -w GOPRIVATE=goa.design/clue
```
