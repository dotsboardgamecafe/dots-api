help: ## You are here! showing all command documenentation.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#===================#
#== Env Variables ==#
#===================#
MIGRATIONS_DIR ?= ./resources/migrations
DATABASE_STRING ?= postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

#========================#
#== DATABASE MIGRATIONS ==#
#========================#
migrate-create: ## Create Database Migration
migrate-create:
	@echo 'Creating migrations file for ${name}'
	migrate create -ext sql -dir ${MIGRATIONS_DIR} -seq ${name}

migrate-up: ## Run Database Migrations
migrate-up:
	@echo 'Running migrations'
	migrate -database ${DATABASE_STRING}  -path ${MIGRATIONS_DIR} -verbose up

migrate-down: ## Drop Database Migrations
migrate-down:
	@echo 'Dropping migrations'
	migrate -database ${DATABASE_STRING}  -path ${MIGRATIONS_DIR} -verbose down -all

migrate-down-force: ## Drop Database Migrations
migrate-down-force:
	@echo 'Dropping migrations'
	migrate -database ${DATABASE_STRING}  -path ${MIGRATIONS_DIR} -verbose force ${version}

migrate: ## Refresh Database Migrations
migrate:
	@echo 'Reseting Database Migrations'
	make migrate-down
	make migrate-up

run: ## Running Binary
run:
	go run main.go
