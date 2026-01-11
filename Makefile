# Build and Manage Q8 Agent

.PHONY: build docker-build docker-pull docker-tag docker-push up down restart logs dev clean help init-config

# Go parameters
GO=$(shell if command -v go >/dev/null 2>&1; then command -v go; elif [ -f /usr/local/go/bin/go ]; then echo /usr/local/go/bin/go; else echo go; fi)
BINARY_NAME=q8-agent
MAIN_PATH=./cmd/agent

# Docker parameters
REGISTRY=91.99.168.0:5000
IMAGE_NAME=q8-agent
TAG=latest

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the Go binary locally
	$(GO) build -o tmp/$(BINARY_NAME) $(MAIN_PATH)

docker-build: ## Build the Docker image locally
	docker build -t $(IMAGE_NAME):$(TAG) .

docker-tag: ## Tag the image for the registry
	docker tag $(IMAGE_NAME):$(TAG) $(REGISTRY)/$(IMAGE_NAME):$(TAG)

docker-push: docker-build docker-tag ## Build, tag and push the image to the registry
	docker push $(REGISTRY)/$(IMAGE_NAME):$(TAG)

docker-pull: ## Pull the latest image from the registry
	docker pull $(REGISTRY)/$(IMAGE_NAME):$(TAG)

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
	$(GO) clean

init-config: ## Create a template .env file with a unique 32-char token
	@if [ ! -f .env ]; then \
		TOKEN=$$(openssl rand -hex 16); \
		echo "Q8_AGENT_PORT=8080" > .env; \
		echo "Q8_AGENT_ADMIN_TOKEN=$$TOKEN" >> .env; \
		echo "Q8_TENANTS_ROOT=/opt/tenants" >> .env; \
		echo "✅ .env file created with unique token: $$TOKEN"; \
	else \
		echo "⚠️  .env file already exists"; \
	fi
