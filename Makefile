include db.mk

APP_NAME := FoodTinder

.PHONY: up down restart ps build run test lint

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose down
	docker compose up -d

ps:
	docker compose ps

build:
	go build -o $(APP_NAME) ./...

run: up
	go run ./cmd/...

test:
	@echo "Running tests..."
	go test ./...

lint:
	golangci-lint run ./...


## dont forget adding swagger