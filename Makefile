.PHONY: tidy build run gen-docs test dev-build dev-up dev-down dev-down-v prod-build prod-up prod-down db-up db-down db-down-v

tidy:
	go mod tidy

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

gen-docs:
	go run ./cmd/gen-docs

test:
	go test ./... -v -cover

dev-build:
	docker compose -f docker-compose.dev.yaml --env-file .env.dev build

dev-up:
	docker compose -f docker-compose.dev.yaml --env-file .env.dev up -d

dev-down:
	docker compose -f docker-compose.dev.yaml --env-file .env.dev down

dev-down-v:
	docker compose -f docker-compose.dev.yaml --env-file .env.dev down -v

prod-build:
	docker compose -f docker-compose.prod.yaml --env-file .env.prod build

prod-up:
	docker compose -f docker-compose.prod.yaml --env-file .env.prod up -d

prod-down:
	docker compose -f docker-compose.prod.yaml --env-file .env.prod down

db-up:
	docker compose -f docker-compose.db.yaml --env-file .env.prod up -d

db-down:
	docker compose -f docker-compose.db.yaml --env-file .env.prod down

db-down-v:
	docker compose -f docker-compose.db.yaml --env-file .env.prod down -v
