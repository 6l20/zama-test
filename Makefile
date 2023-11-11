# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Pattern to watch (*.go)
WATCH_PATTERN=*.go



# Binary names
SERVER_BINARY_NAME=server
CLIENT_BINARY_NAME=client

# Docker parameters
DOCKER_COMPOSE_CMD=docker-compose

all: test build
build: 
	$(GOBUILD) -o $(SERVER_BINARY_NAME) -v ./cmd/server
	$(GOBUILD) -o $(CLIENT_BINARY_NAME) -v ./cmd/client
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(SERVER_BINARY_NAME)
	rm -f $(CLIENT_BINARY_NAME)
run-server: build
	./$(SERVER_BINARY_NAME)
run-client: build
	./$(CLIENT_BINARY_NAME)

# Docker tasks
docker-build:
	$(DOCKER_COMPOSE_CMD) build
docker-up:
	$(DOCKER_COMPOSE_CMD) up
docker-down:
	$(DOCKER_COMPOSE_CMD) down
docker-clean: docker-down
	docker rmi $$(docker images -q --filter "dangling=true") -f

.PHONY: all build test clean run-server run-client docker-build docker-up docker-down docker-clean
