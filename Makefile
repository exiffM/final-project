BIN_MONITOR := "./bin/monitor/monitor"
BIN_CLIENT := "./bin/client/client"
DOCKER_MONITOR_IMG="monitor:develop"
DOCKER_CLIENT_IMG="client:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build-monitor:
	go build -v -o $(BIN_MONITOR) -ldflags "$(LDFLAGS)" ./cmd/monitor

build-client:
	go build -v -o $(BIN_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client

run-monitor: build-monitor
	$(BIN_MONITOR) --config=./configs/config.yml

run-client: build-client
	$(BIN_CLIENT)

build-img-monitor:
	docker build \
	--build-arg=LDFLAGS="$(LDFLAGS)" \
	-t $(DOCKER_MONITOR_IMG) \
	-f build/monitor/Dockerfile .

build-img-client:
	docker build \
	--build-arg=LDFLAGS="$(LDFLAGS)" \
	-t $(DOCKER_CLIENT_IMG) \
	-f build/client/Dockerfile .

run-img-monitor: build-img-monitor
	docker run -p 50051:50051 $(DOCKER_MONITOR_IMG)

run-img-client: build-img-client
	docker run $(DOCKER_CLIENT_IMG)	

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
	go test -v -race -cover -timeout=3m30sm -count=10 ./...

docker-start-components:
	docker compose up --build;

docker-stop:
	docker compose down;

integration-test:
	go test -tags integration -race -v -count=10 -timeout=3m30s ./tests/integration

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint -v run ./...

.PHONY: build run test lint