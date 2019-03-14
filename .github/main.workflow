workflow "Main" {
  on = "push"
  resolves = [ "Tests" ]
}

action "Tests" {
  uses = "docker://golang:1.11"
  args = [ "go", "test", "-race", "./..." ]
}
