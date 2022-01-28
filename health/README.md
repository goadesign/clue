# health: Healthy Services

[![Build Status](https://github.com/goadesign/clue/workflows/CI/badge.svg?branch=main&event=push)](https://github.com/goadesign/clue/actions?query=branch%3Amain+event%3Apush)

## Overview

Package `health` provides a standard health check HTTP endpoint typically served
under the `/healthz` and/or `/livez` paths.

The handler implementation iterates through a given list of service dependencies
and respond with HTTP status `200 OK` if all dependencies are healthy,
`503 Service Unavailable` otherwise. the response body lists each dependency
with its status.

Healthy service example (using the [httpie](https://httpie.org/) command line utility):

```bash
http http://localhost:8083/livez
HTTP/1.1 200 OK
Content-Length: 109
Content-Type: text/plain; charset=utf-8
Date: Mon, 17 Jan 2022 23:23:12 GMT

{
    "status": {
        "ClickHouse": "OK",
        "poller": "OK"
    },
    "uptime": 20,
    "version": "91bb64a8103b494d0eac680f8e929e74882eea5f"
}
```

Unhealthy service:

```bash
http http://localhost:8083/livez
HTTP/1.1 503 Service Unavailable
Content-Length: 113
Content-Type: text/plain; charset=utf-8
Date: Mon, 17 Jan 2022 23:23:20 GMT

{
    "status": {
        "ClickHouse": "OK",
        "poller": "NOT OK"
    },
    "uptime": 20,
    "version": "91bb64a8103b494d0eac680f8e929e74882eea5f"
}
```

## Usage

```go
package main

import (
        "context"
        "database/sql"

       "github.com/goadesign/clue/health"
       "github.com/goadesign/clue/log"
       goahttp "goa.design/goa/v3/http"

        "github.com/repo/services/svc/clients/storage"
       	httpsvrgen "github.com/repo/services/svc/gen/http/svc/server"
       	svcgen "github.com/repo/services/svc/gen/svc"
)

func main() {
        // Initialize the log context
	ctx := log.With(log.Context(context.Background()), "svc", svcgen.ServiceName)

        // Create service clients used by this service
        // The client object must implement the `health.Pinger` interface
	// dsn := ...
	con, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Error(ctx, "could not connect to clickhouse", "err", err.Error())
	}
        stc := storage.New(con)

        // Create the service (user code)
        svc := svc.New(ctx, stc)
        // Wrap the service with Goa endpoints
        endpoints := svcgen.NewEndpoints(svc)

        // Create HTTP server
        mux := goahttp.NewMuxer()
        httpsvr := httpsvrgen.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
        httpsvrgen.Mount(mux, httpsvr)

        // ** Mount health check handler **
	check := health.Handler(health.NewChecker(stc))
	mux.Handle("GET", "/healthz", check)
	mux.Handle("GET", "/livez", check)

        // ... start HTTP server
}
```

Creating an health check HTTP handler is as simple as:

  1. instantiating a health checker using the `NewChecker` function
  2. wrapping it in a HTTP handler using the `Handler` function

The `NewChecker` function accepts a list of dependencies to be checked that must
implement by the `Pinger` interface. 

```go
// Pinger makes it possible to ping a downstream service.
Pinger interface {
	// Name of remote service.
	Name() string
	// Ping the remote service, return a non nil error if the
	// service is not available.
	Ping(context.Context) error
}
```

## Implementing the Pinger Interface

### For Downstream Microservices

The `NewPinger` function instantiates a `Pinger` for a service equipped with a
`/livez` health check endpoint (e.g. a service exposing the handler created by
this package `Handler` function).

### For SQL Databases (e.g. PostgreSQL, ClickHouse)

The stdlib `sql.DB` type provides a `PingContext` method that can be used to
ping a database. Implementing `Pinger` thus consists of adding the following two
methods to the client struct:

```go
// SQL database client used by service.
type client struct {
	db *sql.DB
}

// Ping implements the `health.Pinger` interface.
func (c *client) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// Name implements the `health.Pinger` interface.
func (c *client) Name() string {
	return "PostgreSQL" // ClickHouse, MySQL, etc.
}
```



