.EXPORT_ALL_VARIABLES:
.PHONY:
.SILENT:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := start

include .env

dev:
	air

start:
	GIN_MODE=release go run cmd/main.go

art:
	artillery run test/art.yml

#testing variables
export TEST_DB_URI=mongo://$(MONGODB_USER):$(MONGODB_PASSWORD)@$(MONGODB_HOST):$(MONGODB_PORT)
export TEST_DB_NAME=test
export TEST_CONTAINER_NAME=test_db

tests:
	GIN_MODE=release go test -v ./test/... -coverpkg=./... -cover -coverprofile=coverage.out -failfast -count=1

cover:
	go tool cover -html=coverage.out

lint:
	golangci-lint run

doc:
	swag init -g cmd/main.go

env:
	echo $$TEST_DB_URI