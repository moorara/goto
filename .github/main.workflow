workflow "Main" {
  on = "push"
  resolves = [ "Lint", "Test" ]
}

action "Lint" {
  uses = "docker://golangci/golangci-lint:latest"
  runs = [ "golangci-lint", "run" ]
  args = [ "--deadline", "5m", "--new" ]
}

action "Test" {
  uses = "docker://golang:1.12"
  runs = [ "go", "test" ]
  args = [ "-race", "./..." ]
}

/* action "Cover" {
  uses = "./.github/action-cover"
  secrets = [ "CODECOV_TOKEN" ]
  args = [ "./..." ]
} */
