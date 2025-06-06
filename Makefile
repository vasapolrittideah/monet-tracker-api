# === General Config ===
DOCKER_COMPOSE = docker compose
COMPOSE_FILE = docker-compose.yml
PROTO_DIR=protobuf
OUT_DIR=generated/protobuf
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

# === Docker Compose Targets ===
.PHONY: up down restart logs build

up:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d

down:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down

restart: down up

logs:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f

build:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) build

# === Bazel Targets ===
.PHONY: bazel-gazelle bazel-build bazel-tidy bazel-test

bazel-gazelle:
	bazel run //:gazelle

bazel-build:
	bazel build //...

bazel-tidy:
	bazel mod tidy

bazel-test:
	bazel test //...

bazel-clean:
	bazel clean

# === Combined Commands ===
.PHONY: init all

init: bazel-gazelle bazel-tidy

all: bazel-gazelle bazel-tidy bazel-build bazel-test

# === Protobuf Targets ===
.PHONY: proto
proto:
	protoc \
		-I=./protobuf \
		--go_out=./protogen \
		--go_opt=paths=source_relative \
		--go-grpc_out=./protogen \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=logtostderr=true:./protogen \
		./protobuf/*.proto
	