# makeを打った時のコマンド
.DEFAULT_GOAL := help

# 環境変数
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
IMAGE_TAG=$(git rev-parse HEAD)
ECR_REGISTRY=${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com

.PHONY: build
build: ## Build docker image to deploy
	docker image build \
		-t ${ECR_REGISTRY}/point-app-backend:latest \
		-t ${ECR_REGISTRY}/point-app-backend:${IMAGE_TAG}  \
		--target deploy ./

.PHONY: build-up
build-up: ## Build docker image and up container
	docker compose up -d --build

.PHONY: push
push: ## push to ECR
	aws ecr --region ap-northeast-1 get-login-password | docker login --username AWS --password-stdin https://${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/point-app-backend
	docker image push -a ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/point-app-backend

.PHONY: in
in: ## Appのコンテナに入る（ホスト）
	docker compose exec app sh

.PHONY: up
up: ## Do docker compose up with hot reload（ホスト）
	docker compose up -d

.PHONY: down
down: ## Do docker compose down（ホスト）
	@docker compose down

.PHONY: format
log: ## Tail docker compose logs（ホスト）
	@docker compose logs app -f

.PHONY: ps
ps: ## Check container status（ホスト）
	docker compose ps

.PHONY: rsa 
rsa: down build-up ## 全てのコンテナを削除して、ビルドして、起動

.PHONY: dry-migrate
dry-migrate: ## Try migration（マイグレーション時に発行されるDDL確認）
	mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} -P ${DB_PORT} ${DB_NAME} --dry-run < ./_tools/mysql/schema.sql

.PHONY: migrate
migrate:  ## Execute migration（コンテナ）
	mysqldef -u ${DB_USER} -p ${DB_PASSWORD} -h ${DB_HOST} -P ${DB_PORT} ${DB_NAME} < ./_tools/mysql/schema.sql

.PHONY: seed
seed: ## データ挿入（コンテナ）
	mysql ${DB_NAME} -h ${DB_HOST} -u ${DB_USER} -p${DB_PASSWORD} < ./_tools/mysql/seed.sql 

.PHONY: rdm
rdm: ## 送信メールを見る
	@if [ ${CONTAINER_ENV} ]; then \
		curl -v http://aws:4566/_localstack/ses/ | jq . | tail -n 18 | head -n 16; \
	else \
		curl -v http://localhost:4566/_localstack/ses/ | jq . | tail -n 18 | head -n 16; \
	fi

.PHONY: create-key
create-key: ## JWTに必要なkeyを作成する
	openssl genrsa 4096 > ./auth/certificate/secret.pem
	openssl rsa -pubout < ./auth/certificate/secret.pem > ./auth/certificate/public.pem

.PHONY: format
format: ## フォーマット
	# フォーマット
	@if [ ${CONTAINER_ENV} ]; then \
		gofmt -l -s -w .; \
		goimports -w -l .; \
	else \
		docker compose exec app gofmt -l -s -w .; \
		docker compose exec app goimports -w -l .; \
	fi


.PHONY: lint
lint: format ## リンター(golangci-lint)
	# リンター
	@if [ ${CONTAINER_ENV} ]; then \
		golangci-lint run; \
	else \
		docker compose exec app golangci-lint run; \
	fi


.PHONY: moq
moq: ## mock作成(コンテナ内)
	# サービスのモック生成中
	@docker compose exec app moq -fmt goimports -out ./handler/moq_test.go ./handler \
					RegisterUserService \
					RegisterTemporaryUserService \
					SigninService \
					GetUsersService \
					UpdatePasswordService \
					UpdateAccountService \
					ResetPasswordService \
					SendPointService \
					SignoutService \
					GetAccountService \
					UpdateTemporaryEmailService \
					GetNotificationService \
					GetNotificationsService \
					GetUncheckedNotificationCountService

	# リポジトリのモック生成中
	@docker compose exec app moq -fmt goimports -out ./service/moq_test.go -skip-ensure -pkg service ./domain \
					Cache \
					TokenGenerator \
					UserRepo \
					PointRepo \
					NotificationRepo
	@docker compose exec app moq -fmt goimports -out ./service/repogitory_moq_test.go -skip-ensure -pkg service ./repository Beginner Preparer Execer Queryer Transacter

.PHONY: test
test: ## テスト
	# テスト実行中
	@if [ ${CONTAINER_ENV} ]; then \
		go test -cover -race -shuffle=on ./...; \
	else \
		docker compose exec app go test -cover -race -shuffle=on ./...; \
	fi

.PHONY: mc
mc: ## make coverage カバレッジファイル作成（ホスト側）
	# テスト実行中
	@docker compose exec app go test -cover ./... -coverprofile=cover.out
	# HTMLに変換中
	@docker compose exec app go tool cover -html=cover.out -o tmp/cover.html
	@docker compose exec app rm cover.out

.PHONY: wc
wc: mc ## watch coverage カバレッジを見る（ホスト側）
	# ブラウザで表示
	@open ./tmp/cover.html

.PHONY: wire
wire: ## DIファイル生成
	@wire ./router

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
