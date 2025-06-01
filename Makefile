.PHONY: dev env lint test fmt

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

## run unit tests
test:
	go test ./...
	micromamba run -n gohomeai pytest ml/tests

## auto-format everything
fmt:
	go fmt ./...
	micromamba run -n gohomeai black ml/