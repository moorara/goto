[![Build Status][travisci-image]][travisci-url]

# goto
A collection of handy data structures, algorithms, tools, etc. for [go](https://golang.org) applications.

## Packages

| Package                | Description                                    |
|------------------------|------------------------------------------------|
| [dt](./dt)             | Data type definitions                          |
| [ds](./ds)             | Miscellaneous data structures                  |
| [sort](./sort)         | Common sorting algorithms                      |
| [graphviz](./graphviz) | Generating graphs in **Graphviz DOT** language |
| [config](./config)     | Zero-configuration config management!          |
| [io](./io)             | I/O helper functions                           |
| [math](./math)         | Math helper functions                          |
| [util](./util)         | Miscellaneous helper functions                 |

## Running Tests

| Command          | Purpose                                           |
|------------------|---------------------------------------------------|
| `make test`      | Running unit tests                                |
| `make benchmark` | Running benchmarks                                |
| `make coverage`  | Running unit tests and generating coverage report |


[travisci-url]: https://travis-ci.org/moorara/goto
[travisci-image]: https://travis-ci.org/moorara/goto.svg?branch=master
