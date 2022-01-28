# clue: Microservice Instrumentation

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)
[![codecov](https://codecov.io/gh/goadesign/clue/branch/main/graph/badge.svg?token=HVP4WT1PS6)](https://codecov.io/gh/goadesign/clue)

## Overview

This repository contains microservice instrumentation packages covering the
following topics:

* Logging: the [log](log/) package provides a context-based logging API that
  intelligently selects what to log.
* Metrics: the [instrument](instrument/) package makes it possible for Goa
  services to expose a Prometheus compatible `/metrics` HTTP endpoint.
* Health checks: the [health](health/) package provides a simple way for
  services to expose a health check endpoint.
* Dependency mocks: the [mock](mock/) package provides a way to mock
  downstream dependencies for testing.
* Tracing: the [trace](trace/) package conforms to the
  [OpenTelemetry](https://opentelemetry.io/) specification to trace requests.

Consult the package-specific READMEs for more information.

## Example

The repository contains a [fully functional example](example/weather)
comprised of three instrumented Goa microservices.