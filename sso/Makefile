run: 
	go run cmd/sso/main.go --config=./config/local.yaml
build:
# такой бинарник используеться для деплоя
	go build -o grpc-auth ./cmd/sso
# Запуск миграция
# go run ./cmd/migrator/ --storage-path=./storage/sso.db --migrations-path=./migrations 

.DEFAULT_GOAL := run
.PHONY: run