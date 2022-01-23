# micro: Microservice Instrumentation

[![Build Status](https://github.com/crossnokaye/micro/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/crossnokaye/micro/actions?query=branch%3Amain+event%3Apush)
[![codecov](https://codecov.io/gh/crossnokaye/micro/branch/main/graph/badge.svg?token=HVP4WT1PS6)](https://codecov.io/gh/crossnokaye/micro)

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
  [OpenTelemetry](https://opentelemetry.io/) specification and implements a
  [Goa](https://goa.design) compatible library to trace requests.

Consult the package-specific READMEs for more information.

## Example

The repository contains a [fully functional example](example/weather)
comprised of three instrumented Goa microservices.

## Importing Private Repositories

Make sure to run the following commands before importing any module hosted on
the CrossnoKaye GitHub org private repositories (e.g. this one):

```
git config --global url.git@github.com:.insteadOf https://github.com/
go env -w GOPRIVATE=github.com/crossnokaye/*
```
