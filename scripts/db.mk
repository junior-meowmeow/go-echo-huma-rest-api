ifneq (,$(wildcard config/.env))
include config/.env
else
$(warning config/.env does not exist)
endif

MIGRATION_DIR := db/migrations

DB_HOST ?= $(POSTGRES_HOST)
DB_PORT ?= $(POSTGRES_PORT)
DB_NAME ?= $(POSTGRES_DB)
DB_USER ?= $(POSTGRES_USER)
DB_PASS ?= $(POSTGRES_PASS)

DB_DSN := "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

GOOSE := goose -dir $(MIGRATION_DIR) postgres $(DB_DSN)

.PHONY: migrate-status migrate-up migrate-down migrate-create

migrate-status:
	$(GOOSE) status

migrate-up:
	$(GOOSE) up

migrate-down:
	$(GOOSE) down

migrate-create:
	@$(if $(strip $(name)),,$(error migration name is required. usage: make migrate-create name=new_table))
	goose -dir $(MIGRATION_DIR) create $(name) sql

.PHONY: sqlc

sqlc-gen:
	sqlc generate
