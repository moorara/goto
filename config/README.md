# config

This is a very minimal utility for reading configuration values from either
**command-line flags**, **environment variables**, or **configuration files**.

This library does not use `flag` package for parsing flags, so you can still parse your flags separately.

## Quick Start

```go
package main

import (
  "fmt"
  "net/url"
  "time"

  "github.com/moorara/goto/config"
)

Config := struct {
  Enabled   bool
  LogLevel  string
  Timeout   time.Duration
  Address   url.URL
  Endpoints []string
} {
  Enabled:  true,   // default
  LogLevel: "info", // default
}

func main() {
  config.Pick(&Config)
  fmt.Printf("%+v\n", Config)
}
```

You can run your application with `-help` or `--help` flag to see all options for passing the configuration values.

The precendence of sources for values is as follows:

  1. command-line flags
  2. environment variables
  3. configuration files
  4. default values (set when creating `spec`)

You can pass the configuration values using **flags** using any of the syntaxes below:

```bash
main  -enabled  -log.level info  -timeout 30s  -address http://localhost:8080  -endpoints url1,url2,url3
main  -enabled  -log.level=info  -timeout=30s  -address=http://localhost:8080  -endpoints=url1,url2,url3
main --enabled --log.level info --timeout 30s --address http://localhost:8080 --endpoints url1,url2,url3
main --enabled --log.level=info --timeout=30s --address=http://localhost:8080 --endpoints=url1,url2,url3
```

You can pass the configuration values using **environment variables** as follows:

```bash
export ENABLED=true
export LOG_LEVEL=info
export TIMEOUT=30s
export ADDRESS=http://localhost:8080
export ENDPOINTS=url1,url2,url3
```

You can also write the configuration values in **files**
and set the paths to the files using environment variables:

```bash
export ENABLED_FILE=...
export LOG_LEVEL_FILE=...
export TIMEOUT_FILE=...
export ADDRESS_FILE=...
export ENDPOINTS_FILE=...
```

## Complete Example

```go
package main

import (
  "fmt"
  "net/url"
  "time"

  "github.com/moorara/goto/config"
)

type Config struct {
  unexported         string          // Unexported, will be skipped
  FieldString        string          `flag:"f.string" env:"F_STRING" file:"F_STRING_FILE"`
  FieldBool          bool            `flag:"f.bool" env:"F_BOOL" file:"F_BOOL_FILE"`
  FieldFloat32       float32         `flag:"f.float32" env:"F_FLOAT32" file:"F_FLOAT32_FILE"`
  FieldFloat64       float64         `flag:"f.float64" env:"F_FLOAT64" file:"F_FLOAT64_FILE"`
  FieldInt           int             `flag:"f.int" env:"F_INT" file:"F_INT_FILE"`
  FieldInt8          int8            `flag:"f.int8" env:"F_INT8" file:"F_INT8_FILE"`
  FieldInt16         int16           `flag:"f.int16" env:"F_INT16" file:"F_INT16_FILE"`
  FieldInt32         int32           `flag:"f.int32" env:"F_INT32" file:"F_INT32_FILE"`
  FieldInt64         int64           `flag:"f.int64" env:"F_INT64" file:"F_INT64_FILE"`
  FieldUint          uint            `flag:"f.uint" env:"F_UINT" file:"F_UINT_FILE"`
  FieldUint8         uint8           `flag:"f.uint8" env:"F_UINT8" file:"F_UINT8_FILE"`
  FieldUint16        uint16          `flag:"f.uint16" env:"F_UINT16" file:"F_UINT16_FILE"`
  FieldUint32        uint32          `flag:"f.uint32" env:"F_UINT32" file:"F_UINT32_FILE"`
  FieldUint64        uint64          `flag:"f.uint64" env:"F_UINT64" file:"F_UINT64_FILE"`
  FieldDuration      time.Duration   `flag:"f.duration" env:"F_DURATION" file:"F_DURATION_FILE"`
  FieldURL           url.URL         `flag:"f.url" env:"F_URL" file:"F_URL_FILE"`
  FieldStringArray   []string        `flag:"f.string.array" env:"F_STRING_ARRAY" file:"F_STRING_ARRAY_FILE" sep:","`
  FieldFloat32Array  []float32       `flag:"f.float32.array" env:"F_FLOAT32_ARRAY" file:"F_FLOAT32_ARRAY_FILE" sep:","`
  FieldFloat64Array  []float64       `flag:"f.float64.array" env:"F_FLOAT64_ARRAY" file:"F_FLOAT64_ARRAY_FILE" sep:","`
  FieldIntArray      []int           `flag:"f.int.array" env:"F_INT_ARRAY" file:"F_INT_ARRAY_FILE" sep:","`
  FieldInt8Array     []int8          `flag:"f.int8.array" env:"F_INT8_ARRAY" file:"F_INT8_ARRAY_FILE" sep:","`
  FieldInt16Array    []int16         `flag:"f.int16.array" env:"F_INT16_ARRAY" file:"F_INT16_ARRAY_FILE" sep:","`
  FieldInt32Array    []int32         `flag:"f.int32.array" env:"F_INT32_ARRAY" file:"F_INT32_ARRAY_FILE" sep:","`
  FieldInt64Array    []int64         `flag:"f.int64.array" env:"F_INT64_ARRAY" file:"F_INT64_ARRAY_FILE" sep:","`
  FieldUintArray     []uint          `flag:"f.uint.array" env:"F_UINT_ARRAY" file:"F_UINT_ARRAY_FILE" sep:","`
  FieldUint8Array    []uint8         `flag:"f.uint8.array" env:"F_UINT8_ARRAY" file:"F_UINT8_ARRAY_FILE" sep:","`
  FieldUint16Array   []uint16        `flag:"f.uint16.array" env:"F_UINT16_ARRAY" file:"F_UINT16_ARRAY_FILE" sep:","`
  FieldUint32Array   []uint32        `flag:"f.uint32.array" env:"F_UINT32_ARRAY" file:"F_UINT32_ARRAY_FILE" sep:","`
  FieldUint64Array   []uint64        `flag:"f.uint64.array" env:"F_UINT64_ARRAY" file:"F_UINT64_ARRAY_FILE" sep:","`
  FieldDurationArray []time.Duration `flag:"f.duration.array" env:"F_DURATION_ARRAY" file:"F_DURATION_ARRAY_FILE" sep:","`
  FieldURLArray      []url.URL       `flag:"f.url.array" env:"F_URL_ARRAY" file:"F_URL_ARRAY_FILE" sep:","`
}

func main() {
  c := Config{
    FieldString: "default value",
  }
  config.Pick(&c)
  fmt.Printf("%+v\n", c)
}
```
