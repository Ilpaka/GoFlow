BACKEND_DIR := backend
COMPOSE_DIR := backend/deployments
COMPOSE     := docker compose
PKG         := ./...

.PHONY: help run test test-race cover lint fmt tidy build docker-up docker-down docker-logs docker-build clean migrate

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make <target>\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

run: ## Run app locally (CONFIG_PATH=configs/local.yaml)
	cd $(BACKEND_DIR) && CONFIG_PATH=$${CONFIG_PATH:-configs/local.yaml} go run ./cmd/app

build: ## Build the app binary into bin/app
	cd $(BACKEND_DIR) && go build -trimpath -ldflags="-s -w" -o ../bin/app ./cmd/app

test: ## Run unit tests
	cd $(BACKEND_DIR) && go test $(PKG)

test-race: ## Run tests with the race detector
	cd $(BACKEND_DIR) && go test -race $(PKG)

cover: ## Run tests with coverage profile
	cd $(BACKEND_DIR) && go test -coverprofile=../coverage.out $(PKG)
	@echo "Coverage report: coverage.out"

lint: ## Run golangci-lint
	cd $(BACKEND_DIR) && golangci-lint run

fmt: ## Format Go sources
	cd $(BACKEND_DIR) && gofmt -s -w .

tidy: ## Tidy go.mod / go.sum
	cd $(BACKEND_DIR) && go mod tidy

docker-up: ## Bring up the full compose stack (app + Postgres + Redis + Kafka + observability)
	cd $(COMPOSE_DIR) && $(COMPOSE) up --build -d

docker-down: ## Stop the compose stack
	cd $(COMPOSE_DIR) && $(COMPOSE) down

docker-logs: ## Tail app logs
	cd $(COMPOSE_DIR) && $(COMPOSE) logs -f app

docker-build: ## Build the app image only
	cd $(COMPOSE_DIR) && $(COMPOSE) build app

migrate: ## Migrations run automatically on app startup (RUN_MIGRATIONS_ON_STARTUP=true). This target is a reminder.
	@echo "Migrations are applied at startup when RUN_MIGRATIONS_ON_STARTUP=true."
	@echo "SQL files live in backend/internal/migration/ and are tracked in the schema_migrations table."

clean: ## Remove build artifacts
	rm -rf bin coverage.out coverage.html
