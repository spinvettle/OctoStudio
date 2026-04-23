# Makefile for OctoStudio

.PHONY: pre-commit run web-run docker-build docker-run help

# Run pre-commit hooks on all files
pre-commit:
	pre-commit run --all-files

run:
	go run ./cmd/studio/main.go

web-run:
	pnpm --dir web dev --host 0.0.0.0

docker-build:
	docker build -f dockerfile.yaml -t octostudio-app .

docker-run:
	docker run -d -p 8080:8080 --name my-studio octostudio-app

# Show available targets
help:
	@echo "Available targets:"
	@echo "  pre-commit - Run pre-commit hooks on all files"
	@echo "  run - Start the backend service"
	@echo "  web-run - Start the frontend dev server"
	@echo "  docker-build - Build the Docker image"
	@echo "  docker-run - Run the Docker container"
	@echo "  help - Show this help message"
# llama-server-qwen
# 	./bin/llama-b7974/llama-server \
# 	-m ./storage/models/qwen/Qwen3-0.6B-Q8_0.gguf \
# 	-c 512 \
# 	-t 2 \
# 	--host 0.0.0.0 \
#   	--port 8080 \
