# http

This package provides utilities for HTTP servers.

| Item                           | Description                                                                          |
|--------------------------------|--------------------------------------------------------------------------------------|
| `http.Error`                   | An `error` type capturing context and information about a failed http request.       |
| `http.ResponseWriter`          | An implementation of standard `http.ResponseWriter` for recording status code.       |
| `http.ObservabilityMiddleware` | A middleware providing wrappers for http handlers for logging, metrics, and tracing. |

## Quick Start

You can see an example of using the observability middleware [here](./example).
