workflow "Main" {
  on = "push"
  resolves = [ "Test", "Lint" ]
}

action "Test" {
  uses = "docker://golang:1.12"
  runs = [ "go", "test" ]
  args = [ "-race", "./..." ]
}

action "Lint" {
  uses = "docker://golangci/golangci-lint:latest"
  runs = [ "golangci-lint", "run" ]
  args = [ "--deadline", "5m", "--new" ]
}
