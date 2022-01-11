# micro: Microservice Instrumentation

---
[![Build Status](https://github.com/crossnokaye/micro/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/crossnokaye/micro/actions?query=branch%3Amain+event%3Apush)

## Overview

This repository contains microservice instrumentation packages covering the
following topics:

* Logging: the [log](log/) package provides a context-based logging API that
  intelligently selects what to log.
* Tracing: the [tracing](tracing/) package conforms to the
  [OpenTelemetry](https://opentelemetry.io/) specification and implements a
  [Goa](https://goa.design) compatible library to trace requests.
* Metrics: the [metrics](metrics/) package provides makes it possible for Goa
  services to expose a Prometheus compatible metrics endpoint.

Consult the package-specific READMEs for more information.

TBD:
- [x] Logging
- [ ] Tracing
- [ ] Metrics

## Importing Private Repositories

Make sure to run the following commands before importing any module hosted on
the CrossnoKaye GitHub org private repositories (e.g. this one):

```
git config --global url.git@github.com:.insteadOf https://github.com/
go env -w GOPRIVATE=github.com/crossnokaye/*
```