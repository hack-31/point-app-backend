# ----------------------------------------------
# ビルド用環境
FROM golang:1.19.2-bullseye AS deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# ----------------------------------------------
# 本番環境
FROM debian:bullseye-slim AS deploy

# X509: Certificate Signed by Unknown Authorityエラーを回避する
RUN apt-get update \
 && apt-get install -y --force-yes --no-install-recommends apt-transport-https curl ca-certificates \
 && apt-get clean \
 && apt-get autoremove \
 && rm -rf /var/lib/apt/lists/*

COPY --from=deploy-builder /app/app .

CMD ["./app"]

# ----------------------------------------------
# 開発環境
FROM golang:1.19.2-alpine3.16 AS dev

WORKDIR /app

RUN apk update && apk add alpine-sdk && apk add jq

RUN go install github.com/cosmtrek/air@latest 
RUN go install github.com/k0kubun/sqldef/cmd/mysqldef@latest
