# log

This package provides a logger for your go applications.
Default output format is `log.JSON` and default log level is `log.InfoLevel`. 

## Quick Start

You can use the **global/singelton** logger as follows:

```go
package main

import "github.com/moorara/goto/log"

func main() {
  log.SetOptions(log.Options{
    Environment: "prod",
    Region:      "us-east-1",
    Component:   "auth-service",
  })

  log.Error(
    "message", "Hello, World!",
    "error", errors.New("too late!"),
  )
}
```

Output:

```json
{"caller":"main.go:12","component":"auth-service","environment":"prod","error":"too late!","level":"error","message":"Hello, World!","region":"us-east-1","timestamp":"2019-02-12T17:59:33.973456Z"}
```

Or you can create a new instance logger as follows:

```go
package main

import "github.com/moorara/goto/log"

func main() {
  logger := log.NewLogger(log.Options{
    Format:      log.JSON,
    Level:       "debug",
    Name:        "handler",
    Environment: "stage",
    Region:      "us-east-1",
    Component:   "auth-service",
  })

  logger.Debug(
    "message", "Hello, World!",
    "context", map[string]interface{}{
      "retries": 4,
    },
  )
}
```

Output:

```json
{"caller":"main.go:15","component":"auth-service","context":{"retries":4},"environment":"prod","level":"debug","logger":"instance","message":"Hello, World!","region":"us-east-1","timestamp":"2019-02-12T17:59:33.973595Z"}
```
