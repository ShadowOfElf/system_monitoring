BIN := "./bin/system_monitor"
DOCKER_IMG="monitor:develop"
GOOS = linux
GOARCH = amd64
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	GOOS=darwin GOARCH=arm64 go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd

build-linux:
	GOOS=linux GOARCH=amd64 go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd

build-win:
	go build -v -o $(BIN).exe ./cmd

run:
	$(BIN) --config=./bin/test.toml

build-img:
	docker compose -f ./docker_compose/docker-compose.yml build

run-img:
	docker compose -f ./docker_compose/docker-compose.yml up -d

lint:
	golangci-lint run ./...

generate:
	rm -rf ./pkg
	mkdir -p ./pkg
	protoc --go_out=./pkg --go-grpc_out=./pkg/ --proto_path=./protobuf ./protobuf/*.proto

test:
	go test -v -race ./internal/...

