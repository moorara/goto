#!/usr/bin/env bash

set -euo pipefail

go test -cover -covermode=atomic -coverprofile=c.out ./...
go tool cover -html=c.out -o cover.html
bash <(curl -s https://codecov.io/bash) -f c.out
