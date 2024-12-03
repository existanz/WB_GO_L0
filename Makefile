include .env

all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

n ?= 1
send:
	@go run cmd/publish/main.go -n=$(n)

infra:
	docker compose up --build

test:
	@echo "Testing..."
	@go test ./... -v

clean:
	@echo "Cleaning..."
	@rm -f main

migrate-up:
	goose -dir ./migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_DATABASE} search_path=${DB_SCHEMA} sslmode=disable" up

migrate-down:
	goose -dir ./migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_DATABASE} search_path=${DB_SCHEMA} sslmode=disable" down