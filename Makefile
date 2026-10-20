GO ?= go

.PHONY: setup run test vet build fmt infra-up infra-down migrate-up migrate-down migrate-version

setup: migrate-up test vet build
	$(GO) run ./cmd/checksetup

run:
	$(GO) run ./cmd/server

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

build:
	$(GO) build -o bin/server ./cmd/server

fmt:
	$(GO) fmt ./...

lint:
	go tool golangci-lint run --fix -c .golangci.yml ./...

infra-up:
	docker compose up -d --wait postgres redis

infra-down:
	docker compose down

migrate-up: infra-up
	docker compose run --rm migrate up

# Destructive: rolls back the latest migration, dropping attempts at version 1.
migrate-down: infra-up
	docker compose run --rm migrate down 1

migrate-version: infra-up
	docker compose run --rm migrate version
