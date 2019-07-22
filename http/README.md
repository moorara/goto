# http

This package provides utilities for HTTP servers and clients.

| Item                    | Description                                                                                       |
|-------------------------|---------------------------------------------------------------------------------------------------|
| `http.Error`            | An `error` type capturing context and information about a failed http request.                    |
| `http.ResponseWriter`   | An implementation of standard `http.ResponseWriter` for recording status code.                    |
| `http.ClientMiddleware` | A client-side middleware providing wrappers for logging, metrics, tracing, etc.                   |
| `http.ServerMiddleware` | A server-side middleware providing wrappers for http handlers for logging, metrics, tracing, etc. |

## Quick Start

You can see an example of using the middleware [here](./example).
