# micro: Microservice Instrumentation

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