.PHONY: dev env lint api-ml mlflow-ui

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

## run ml api (dev)
api-ml:
	fastapi dev ./ml/main.py --port 9999

## run mlflow-ui
mlflow-ui:
	mlflow ui --port 5111