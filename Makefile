# Makefile for OctoStudio

.PHONY: pre-commit help

# Run pre-commit hooks on all files
pre-commit:
	pre-commit run --all-files
run:
	go run ./cmd/main.go
# Show available targets
help:
	@echo "Available targets:"
	@echo "  pre-commit - Run pre-commit hooks on all files"
	@echo "  help - Show this help message"
# llama-server-qwen
# 	./bin/llama-b7974/llama-server \
# 	-m ./storage/models/qwen/Qwen3-0.6B-Q8_0.gguf \
# 	-c 512 \
# 	-t 2 \
# 	--host 0.0.0.0 \
#   	--port 8080 \
