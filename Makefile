# Variables for Docker image names and tags
DOCKER_SERVER_IMAGE ?= word-of-wisdom-server
DOCKER_CLIENT_IMAGE ?= word-of-wisdom-client
DOCKER_TAG ?= latest

.PHONY: help
all: help
help: Makefile
	@echo
	@echo "Choose a command to run in "$(APP_NAME)":"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

.PHONY: docker-build-server
docker-build-server: ## Build the Docker image for the server
	docker build -t $(DOCKER_SERVER_IMAGE):$(DOCKER_TAG) -f build/server.Dockerfile .

.PHONY: docker-build-client
docker-build-client: ## Build the Docker image for the client
	docker build -t $(DOCKER_CLIENT_IMAGE):$(DOCKER_TAG) -f build/client.Dockerfile .

.PHONY: docker-build-all
docker-build-all: ## Build both server and client Docker images
	docker-build-server docker-build-client

.PHONY: up
up: ## Start server
	docker-compose up --build

.PHONY: down
down: ## Stop and remove containers, networks
	@docker-compose down
