.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

build: ## Build docker image to deploy
	AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
	IMAGE_TAG=$(git rev-parse HEAD)
	ECR_REGISTRY=${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com
	docker image build \
		-t ${ECR_REGISTRY}/point-app-backend:latest \
		-t ${ECR_REGISTRY}/point-app-backend:${IMAGE_TAG}  \
		--target deploy ./

build-l: ## Build docker image to local development
	docker build --no-cache --target deploy ./

build-up: ## Build docker image and up container
	docker compose up -d --build

push: ## push to ECR
	AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
	aws ecr --region ap-northeast-1 get-login-password | docker login --username AWS --password-stdin https://${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/point-app-backend
	docker image push -a ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/point-app-backend

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
	mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} -P ${DB_PORT} ${DB_NAME} --dry-run < ./_tools/mysql/schema.sql

migrate:  ## Execute migration
	mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} -P ${DB_PORT} ${DB_NAME} < ./_tools/mysql/schema.sql

seed: ## seed data to db
	mysql ${DB_NAME} -h ${DB_HOST} -u ${DB_USER} -p${DB_PASSWORD} < ./_tools/mysql/seed.sql 

read-mail-h: ## 送信メールを見る(ホスト側)
	curl -v http://localhost:4566/_localstack/ses/ | jq . | tail -n 18 | head -n 16

read-mail-c: ## 送信メールを見る(コンテナ側)
	curl -v http://aws:4566/_localstack/ses/ | jq . | tail -n 18 | head -n 16

create-key: ## JWTに必要なkeyを作成する
	openssl genrsa 4096 > ./auth/certificate/secret.pem
	openssl rsa -pubout < ./auth/certificate/secret.pem > ./auth/certificate/public.pem

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
