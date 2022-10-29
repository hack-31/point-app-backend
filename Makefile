.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build: ## Build docker image to deploy
	docker build -t budougumi0617/gotodo:${DOCKER_TAG} \
		--target deploy ./

build-local: ## Build docker image to local development
	docker compose build --no-cache

build-up: ## Build docker image and up container
	docker compose up -d --build

serve: ## serve with air 
	docker compose exec app air

up: ## Do docker compose up with hot reload
	docker compose up -d

down: ## Do docker compose down
	docker compose down

logs: ## Tail docker compose logs
	docker compose logs -f

ps: ## Check container status
	docker compose ps

dry-migrate: ## Try migration
	mysqldef -u admin -p password -h db -P 3306 point_app --dry-run < ./_tools/mysql/schema.sql

migrate:  ## Execute migration
	mysqldef -u admin -p password -h db -P 3306 point_app < ./_tools/mysql/schema.sql

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
