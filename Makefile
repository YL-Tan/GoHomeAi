.PHONY: dev env lint mlflow-ui

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

## run mlflow-ui
mlflow-ui:
	@mlflow ui --backend-store-uri ./ml/mlruns --port 5000