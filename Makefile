.DEFAULT_GOAL := help

# Capture extra args (service names)
SERVICES := $(filter-out $@,$(MAKECMDGOALS))

## ---- Core lifecycle ----

up: ## Start containers (optionally specify services)
	docker compose up -d --build $(SERVICES)

down: ## Stop containers
	docker compose down

restart: ## Restart containers
	docker compose down
	docker compose up -d $(SERVICES)

ps: ## Show container status
	docker compose ps $(SERVICES)

logs: ## Follow logs (optionally specify service)
	docker compose logs -f $(SERVICES)

## ---- Build targets ----

build: ## Build images using cache
	docker compose build $(SERVICES)

rebuild: ## Rebuild images without using Docker layer cache
	docker compose build --no-cache $(SERVICES)

up-rebuild: ## Rebuild without cache and start containers
	docker compose up -d --build --no-cache $(SERVICES)

## ---- Cleanup (be careful) ----

clean: ## Stop containers and remove volumes + orphans
	docker compose down -v --remove-orphans

disk-usage: ## Show Docker disk usage
	docker system df
	docker buildx du

prune-images: ## Remove unused images
	docker image prune -a -f

prune-build: ## Remove unused build cache
	docker buildx prune -f

prune-all: ## Remove unused images + build cache (safe)
	docker image prune -a -f
	docker buildx prune -f

## ---- Help ----

help: ## Show this help
	@echo ""
	@echo "Available commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} \
		/^[a-zA-Z_-]+:.*##/ {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
