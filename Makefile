ifeq (,$(wildcard .env))
  $(info .env not found - creating from .env.example)
  $(shell cp .env.example .env)
endif

include .env
export $(shell sed 's/=.*//' .env)

include db.mk

# Ensure gotestsum is installed
GOTESTSUM := $(shell command -v gotestsum 2> /dev/null)
APP_NAME := FoodTinder

.PHONY: up down restart ps build run lint tests-integration tests-unit tests install-gotestsum

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

install-gotestsum:
	@if [ -z "$(GOTESTSUM)" ]; then \
		echo "Installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
	else \
		echo "gotestsum is already installed"; \
	fi

tests-integration: install-gotestsum
	@echo "Running integration tests..."
	@gotestsum  --format=testname ./tests

tests: tests-integration

lint:
	golangci-lint run ./...


## dont forget adding swagger