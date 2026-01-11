# Build and Manage Q8 Agent

.PHONY: build docker-build up down restart logs dev clean help

# Binary name
BINARY_NAME=q8-agent
MAIN_PATH=./cmd/agent

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the Go binary locally
	go build -o tmp/$(BINARY_NAME) $(MAIN_PATH)

docker-build: ## Build the Docker image
	docker compose build

up: ## Start the agent in detached mode
	docker compose up -d

down: ## Stop and remove the agent container
	docker compose down

restart: ## Restart the agent container
	docker compose restart

logs: ## Tail the agent logs
	docker compose logs -f

dev: ## Run with hot-reloading using Air
	air

clean: ## Remove local binary and build artifacts
	rm -rf tmp/
	go clean
