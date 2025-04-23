BIN := "./bin/system_monitor"
DOCKER_IMG="monitor:develop"
GOOS = linux
GOARCH = amd64

build:
	go build -v -o $(BIN) ./cmd

build-linux:
	GOOS=linux GOARCH=amd64 go build -v -o $(BIN) ./cmd

run:
	$(BIN) --config ./configs/test.toml

build-img:
	docker compose -f ./docker_compose/docker-compose.yml build

run-img:
	docker compose -f ./docker_compose/docker-compose.yml up -d
