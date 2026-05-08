# Makefile for OctoStudio

MIGRATIONS_DIR ?= ./migrations
DB_DSN ?= postgres://baishan:baishan@localhost:5432/learn?sslmode=disable
STEP ?= 1

.PHONY: pre-commit run docker-build docker-run migrate-up migrate-down migrate-version migrate-create help

# Run pre-commit hooks on all files
pre-commit:
	pre-commit run --all-files

run:
	go run ./cmd/studio/main.go

docker-build:
	docker build -f dockerfile.yaml -t octostudio-app .

docker-run:
	docker run -d -p 8080:8080 --name my-studio octostudio-app

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down $(STEP)

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" version

migrate-create:
	@test -n "$(NAME)" || (echo "usage: make migrate-create NAME=create_accounts" && exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

# Show available targets
help:
	@echo "Available targets:"
	@echo "  pre-commit - Run pre-commit hooks on all files"
	@echo "  run - Run the server"
	@echo "  docker-build - Build docker image"
	@echo "  docker-run - Run docker container"
	@echo "  migrate-up - Run pending database migrations"
	@echo "  migrate-down - Roll back database migrations, override with STEP=N"
	@echo "  migrate-version - Show current database migration version"
	@echo "  migrate-create - Create migration files, usage: make migrate-create NAME=create_accounts"
	@echo "  help - Show this help message"
# llama-server-qwen
# 	./bin/llama-b7974/llama-server \
# 	-m ./storage/models/qwen/Qwen3-0.6B-Q8_0.gguf \
# 	-c 512 \
# 	-t 2 \
# 	--host 0.0.0.0 \
#   	--port 8080 \
