BIN := "./bin/monitor"
DOCKER_IMG="monitor:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/monitor

run: build
	$(BIN) --config=./configs/config.yml

generate:
	go get \
		google.golang.org/grpc \
		google.golang.org/protobuf/cmd/protoc-gen-go \
    	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	rm -rf ./internal/grpc/pb
	mkdir ./internal/grpc/pb

	protoc --go-grpc_out=./internal/grpc/pb --go_out=./internal/grpc/pb ./api/MonitorService.proto

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint -v run ./...

test: