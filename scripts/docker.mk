COMPOSE := docker compose
ENV_FILE := --env-file ./config/.env

.PHONY: compose-build compose-up compose-down

compose-build:
	$(COMPOSE) build

compose-up:
	$(COMPOSE) ${ENV_FILE} up -d

compose-down:
	$(COMPOSE) ${ENV_FILE} down

DEV_COMPOSE := docker compose -f compose.dev.yaml
DEV_ENV_FILE := --env-file ./config/.env.dev

.PHONY: dev-build dev-up dev-down dev-down-v

dev-build:
	$(DEV_COMPOSE) build

dev-up:
	$(DEV_COMPOSE) ${DEV_ENV_FILE} up -d

dev-down:
	$(DEV_COMPOSE) ${DEV_ENV_FILE} down

dev-down-v:
	$(DEV_COMPOSE) ${DEV_ENV_FILE} down -v