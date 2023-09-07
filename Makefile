BIN_MONITOR := "./bin/monitor/monitor"
BIN_CLIENT := "./bin/client/client"
DOCKER_IMG="monitor:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_MONITOR) -ldflags "$(LDFLAGS)" ./cmd/monitor
	go build -v -o $(BIN_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client

run: build
	$(BIN_MONITOR) --config=./configs/config.yml
	$(BIN_CLIENT)

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version:
	$(BIN_MONITOR) version

generate:
	go get \
		google.golang.org/grpc \
		google.golang.org/protobuf/cmd/protoc-gen-go \
    	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	rm -rf ./internal/grpc/pb
	mkdir ./internal/grpc/pb

	protoc --go-grpc_out=./internal/grpc/pb --go_out=./internal/grpc/pb ./api/MonitorService.proto

test:
	go test -v -race -cover -timeout=1m ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint -v run ./...

.PHONY: build run test lint