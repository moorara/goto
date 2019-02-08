workflow "Main" {
  on = "push"
  resolves = [ "Testing" ]
}

action "Testing" {
  uses = "docker://golang:1.11"
  args = [ "go", "test", "-race", "./..." ]
}
