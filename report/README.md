# report

This package provides an error reporter for your go applications.
An error reporter can report your errors to an error monitoring service ([Rollbar](https://rollbar.com), [Airbrake](https://airbrake.io), etc.)
An error reporting service simplifies creating dashboards, alerting, and debugging.

## Quick Start

You can use the **global/singelton** reporter as follows:

```go
package main

import "github.com/moorara/goto/report"

func main() {
  report.SetOptions(report.RollbarOptions{
    Token:       "rollbar-token",
    Environment: "production",
    CodeVersion: "0.1.0",
    ProjectURL:  "github.com/moorara/repo",
  })

  // Catch panics and report them
  defer report.OnPanic()

  // Report an error
  report.Error(err)
}
```

Or you can create a new instance reporter as follows:

```go
package main

import "github.com/moorara/goto/report"

func main() {
  reporter := report.NewRollbarReporter(report.RollbarOptions{
    Token:       "rollbar-token",
    Environment: "production",
    CodeVersion: "0.1.0",
    ProjectURL:  "github.com/moorara/repo",
  })

  // Catch panics and report them
  defer reporter.OnPanic()

  // Report an error
  reporter.Error(err)
}
```