.PHONY: all build run migrate migrate-down migrate-fresh seed dev

# Build the application
build:
	cd backend && go build -o bin/api ./cmd/api

# Run the application
run: build
	cd backend && ./bin/api

# Run in development mode
dev:
	cd backend && go run ./cmd/api

# Database migrations
migrate:
	cd backend && go run ./cmd/api -migrate

migrate-down:
	cd backend && go run ./cmd/api -migrate-down

migrate-fresh:
	cd backend && go run ./cmd/api -migrate-fresh

# Seed database
seed:
	cd backend && go run ./cmd/api -seed

# Fresh migration + seed
setup: migrate-fresh
	cd backend && go run ./cmd/api -seed

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

# Frontend - Admin Panel
panel-dev:
	cd frontend/panel && npm run dev

panel-build:
	cd frontend/panel && npm run build

panel-install:
	cd frontend/panel && npm install

# Generate JWT Secret
genkey:
	cd backend && go run ./cmd/genkey

genkey-hex:
	cd backend && go run ./cmd/genkey --format hex

genkey-update:
	cd backend && go run ./cmd/genkey --update-env

genkey-bash:
	cd backend && ./scripts/generate_jwt_secret.sh

genkey-bash-update:
	cd backend && ./scripts/generate_jwt_secret.sh --update

# Install all dependencies
deps:
	cd backend && go mod tidy
	cd frontend/panel && npm install
