test:
	@ go test -v -race ./...

benchmark:
	@ go test -run=none -bench=. -benchmem ./...

coverage:
	@ go test -covermode atomic -coverprofile cover.out ./...
	@ go tool cover -html=cover.out -o cover.html


.PHONY: test benchmark coverage
