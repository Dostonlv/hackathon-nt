.PHONY: run_db run test build

APP_NAME=tender-management
DOCKER_COMPOSE=docker-compose.yml

# Database connection parameters
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= tender_db
DB_HOST ?= localhost
DB_PORT ?= 5433
MIGRATION_DIR ?= migrations

# Build connection string
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable



.PHONY: run_db run test build

run_db:
	docker compose up db -d
	docker compose up redis -d

run:
	@echo "Starting services with Docker Compose..."
	@docker compose -f $(DOCKER_COMPOSE) up -d
	@sleep 5 # Allow some time for services to start up
	@echo "Applying database migrations..."
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up
	@echo "Services are running and migrations are applied."


test:
	go test ./...

build:
	go build -o main cmd/server/main.go

stop:
	docker compose down

clean:
	docker compose down -v




# Migration tools
MIGRATE := migrate
MIGRATE_VER := v4.17.0
MIGRATE_PLATFORM := linux
MIGRATE_ARCH := amd64

.PHONY: install-migrate migrate-up migrate-down migrate-create migrate-force migrate-version migrate-status help

# Default target when just running 'make'
.DEFAULT_GOAL := help

# Install golang-migrate if not installed
install-migrate:
	@if ! command -v migrate >/dev/null 2>&1; then \
	echo "Installing golang-migrate..." && \
	curl -L https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VER)/migrate.$(MIGRATE_PLATFORM)-$(MIGRATE_ARCH).tar.gz | tar xvz && \
	sudo mv migrate /usr/local/bin/migrate && \
	echo "golang-migrate installed successfully"; \
	else \
	echo "golang-migrate is already installed"; \
	fi

# Create a new migration file
migrate-create:
	@if [ -z "$(NAME)" ]; then \
	echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
	exit 1; \
	fi
	@echo "Creating migration files for $(NAME)..."
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)

# Run all up migrations
migrate-up:
	@echo "Running up migrations..."
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

# Roll back all migrations
migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down

# Force set migration version
migrate-force:
	@if [ -z "$(VERSION)" ]; then \
	echo "Error: VERSION is required. Usage: make migrate-force VERSION=x"; \
	exit 1; \
	fi
	@echo "Force setting migration version to $(VERSION)..."
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" force $(VERSION)

# Show current migration version
migrate-version:
	@echo "Current migration version:"
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" version

# Show migrations status
migrate-status:
	@echo "Migration status:"
	@migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" status

# Create database if it doesn't exist
create-db:
	@echo "Creating database $(DB_NAME) if it doesn't exist..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$(DB_NAME)'" | grep -q 1 || \
	psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"

# Drop database
drop-db:
	@echo "Dropping database $(DB_NAME)..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"

# Reset database (drop, create, migrate up)
reset-db: drop-db create-db migrate-up
	@echo "Database reset completed successfully"