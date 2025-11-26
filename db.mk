include .env
export $(shell sed 's/=.*//' .env)

DB_URL := postgres://$(DB_USER):$(DB_PASS)@localhost:5432/$(DB_NAME)?sslmode=disable

migrate:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down

migrate-version:
	migrate -path db/migrations -database "$(DB_URL)" version

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=add_users_table"; \
	else \
		migrate create -ext sql -dir db/migrations -seq $(name); \
	fi