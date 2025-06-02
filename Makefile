.PHONY: dev env lint

## create/update Python env
env:
	@bash scripts/bootstrap.sh

## start all services locally
dev: env
	docker compose -f infra/local-compose.yml up --build

## run Go + Python linters
lint:
	golangci-lint run ./...
	micromamba run -n gohomeai ruff ml/
