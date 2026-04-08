-include .env
export

.PHONY: dev dev-detached dev-down migrate-up migrate-down test test-integration lint build setup setup-zitadel install-frontend dev-frontend build-frontend

COMPOSE = docker compose --env-file .env -f infra/docker-compose.yml

dev:
	$(COMPOSE) up --build

dev-detached:
	$(COMPOSE) up --build -d

dev-down:
	$(COMPOSE) down

migrate-up:
	docker exec -e POSTGRES_HOST=postgres -e POSTGRES_PORT=5432 \
		-e POSTGRES_DB=$(POSTGRES_DB) -e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_SSLMODE=disable \
		infra-api-1 go run ./cmd/migrate up

migrate-down:
	docker exec -e POSTGRES_HOST=postgres -e POSTGRES_PORT=5432 \
		-e POSTGRES_DB=$(POSTGRES_DB) -e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_SSLMODE=disable \
		infra-api-1 go run ./cmd/migrate down 1

test:
	cd services/api && go test ./... -v -count=1

test-integration:
	cd services/api && go test ./... -v -count=1 -tags integration

lint:
	cd services/api && go vet ./...

build:
	cd services/api && go build -o ./bin/api ./cmd/api

# Full first-time setup: start stack, wait for health, run migrations, create OIDC app
setup:
	@echo "Starting stack..."
	$(COMPOSE) up --build -d
	@echo "Waiting for all services to be healthy..."
	@until $(COMPOSE) ps --format '{{.Health}}' | grep -v healthy | grep -c . | grep -q '^0$$'; do sleep 5; done 2>/dev/null || sleep 30
	@echo "Running migrations..."
	$(MAKE) migrate-up
	@echo "Setting up Zitadel OIDC..."
	$(MAKE) setup-zitadel
	@echo ""
	@echo "=== Setup complete ==="
	@echo "API:     http://localhost:8081"
	@echo "Zitadel: http://zitadel:8080 (add '127.0.0.1 zitadel' to /etc/hosts)"
	@echo "Health:  curl http://localhost:8081/health"

setup-zitadel:
	@./infra/scripts/setup-zitadel.sh
	@$(COMPOSE) up -d api --force-recreate

# Frontend
install-frontend:
	pnpm install

dev-frontend:
	pnpm turbo dev --filter='./apps/*'

build-frontend:
	pnpm turbo build --filter='./apps/*'
