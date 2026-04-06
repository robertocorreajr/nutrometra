.PHONY: dev test migrate-up migrate-down lint build setup

dev:
	docker compose -f infra/docker-compose.yml up --build

dev-down:
	docker compose -f infra/docker-compose.yml down

migrate-up:
	cd services/api && go run -tags migrate ./cmd/migrate up

migrate-down:
	cd services/api && go run -tags migrate ./cmd/migrate down 1

test:
	cd services/api && go test ./... -v -count=1

test-integration:
	cd services/api && go test ./... -v -count=1 -tags integration

lint:
	cd services/api && go vet ./...

build:
	cd services/api && go build -o ./bin/api ./cmd/api

setup: dev
	@echo "Aguardando Zitadel ficar saudável..."
	@until curl -sf http://localhost:8080/debug/healthz > /dev/null; do sleep 2; done
	@echo "Zitadel pronto. Acesse http://localhost:8080 para configurar o OIDC client."
