.PHONY: tools
tools:
	go get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: zonePB
zonePB:
	cd zonePB && protoc --proto_path=. --go_out=plugins=grpc:. zone.proto

.PHONY: server
server:
	@ cd server && go run .

.PHONY: client
client:
	@ cd client && go run .
