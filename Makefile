SHELL := bash

# Local development defaults. Override on the command line, e.g.
#   make dev PGPORT=5433
PGDATA ?= .postgres
PGPORT ?= 5432
PGUSER ?= postgres
PGPASSWORD ?= postgres
PGDATABASE ?= eventer

DB_HOST ?= localhost
DB_PORT ?= $(PGPORT)
DB_USER ?= $(PGUSER)
DB_PASSWORD ?= $(PGPASSWORD)
DB_NAME ?= $(PGDATABASE)
PORT ?= 8080

export DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME PORT

FRONTEND_DIR := frontend

.DEFAULT_GOAL := help
.PHONY: help dev backend frontend db-init db-start db-stop db-reset migrate \
        test build build-backend build-frontend

help: ## Show the available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

## --- development -------------------------------------------------------------

dev: db-start ## Run Postgres, the API, and the UI together
	$(MAKE) -j2 backend frontend

backend: ## Run the API locally (needs Postgres)
	go run .

frontend: ## Run the Vite dev server
	cd $(FRONTEND_DIR) && { [ -d node_modules ] || npm ci; }
	cd $(FRONTEND_DIR) && npm run dev

## --- local Postgres ----------------------------------------------------------

db-init: ## Initialise and start the local Postgres server
	@[ -f $(PGDATA)/PG_VERSION ] || initdb -D $(PGDATA) -U $(PGUSER) --auth=trust >/dev/null
	@pg_ctl -D $(PGDATA) -o "-p $(PGPORT) -k $(CURDIR)/$(PGDATA)" status >/dev/null 2>&1 || pg_ctl -D $(PGDATA) -o "-p $(PGPORT) -k $(CURDIR)/$(PGDATA)" -l $(PGDATA)/server.log start
	@createdb -h localhost -p $(PGPORT) -U $(PGUSER) $(PGDATABASE) 2>/dev/null || true
	@echo "Postgres ready on localhost:$(PGPORT) (database: $(PGDATABASE))"

db-start: db-init ## Start the local Postgres server

db-stop: ## Stop the local Postgres server
	@pg_ctl -D $(PGDATA) stop -m fast >/dev/null 2>&1 || true

db-reset: ## Stop the server and delete the local Postgres data
	@pg_ctl -D $(PGDATA) stop -m immediate >/dev/null 2>&1 || true
	rm -rf $(PGDATA)

## --- database migrations -----------------------------------------------------

migrate: ## Apply database migrations with goose
	@if [ -d migrations ]; then \
		goose -dir migrations postgres "host=localhost port=$(PGPORT) user=$(PGUSER) password=$(PGPASSWORD) dbname=$(PGDATABASE) sslmode=disable" up; \
	else \
		echo "No migrations directory yet; nothing to do."; \
	fi

## --- tests -------------------------------------------------------------------

test: ## Run the Go tests
	go test ./...

## --- builds ------------------------------------------------------------------

build: build-backend build-frontend ## Build the backend binary and the frontend bundle

build-backend: ## Build the Go binary into bin/
	mkdir -p bin
	go build -o bin/eventer .

build-frontend: ## Build the Vite production bundle
	cd $(FRONTEND_DIR) && { [ -d node_modules ] || npm ci; }
	cd $(FRONTEND_DIR) && npm run build
