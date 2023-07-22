GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	# 这里使用v3
	# go install github.com/micro/micro/v3/cmd/protoc-gen-micro@latest
	go get github.com/asim/go-micro/cmd/protoc-gen-micro/v3
	go install github.com/micro/micro/v3/cmd/protoc-gen-openapi@latest

.PHONY: api
api:
	protoc --openapi_out=. --proto_path=. proto/podapi.proto

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=:. proto/podapi.proto

.PHONY: build
build:
	go build -o pod-api main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t pod-api:latest

.PHONY: run
run:
	go run main.go
