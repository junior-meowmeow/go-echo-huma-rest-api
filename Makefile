.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: run
run:
	go run ./cmd/server

.PHONY: gen-docs
gen-docs:
	go run ./cmd/gen-docs

.PHONY: test
test:
	go test ./... -v

.PHONY: test-cover
test-cover:
	go test ./... -cover

.PHONY: test-clean
test-clean:
	go clean -testcache && go test -v ./... -cover

.PHONY: pre-commit
pre-commit: tidy fmt lint gen-docs test
	@echo "✅ All pre-commit checks passed."

COMPOSE := docker compose
ENV_FILE := --env-file ./config/.env

.PHONY: compose-build
compose-build:
	$(COMPOSE) build

.PHONY: compose-up
compose-up:
	$(COMPOSE) ${ENV_FILE} up -d

.PHONY: compose-down
compose-down:
	$(COMPOSE) ${ENV_FILE} down

DEV_COMPOSE := docker compose -f compose.dev.yaml
DEV_ENV_FILE := --env-file ./config/.env.dev

.PHONY: dev-build
dev-build:
	$(DEV_COMPOSE) build

.PHONY: dev-up
dev-up:
	$(DEV_COMPOSE) ${DEV_ENV_FILE} up -d

.PHONY: dev-down
dev-down:
	$(DEV_COMPOSE) ${DEV_ENV_FILE} down

.PHONY: dev-down-v
dev-down-v:
	$(DEV_COMPOSE) ${DEV_ENV_FILE} down -v
