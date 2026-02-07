# Makefile for OctoStudio

.PHONY: pre-commit help

# Run pre-commit hooks on all files
pre-commit:
	pre-commit run --all-files

# Show available targets
help:
	@echo "Available targets:"
	@echo "  pre-commit - Run pre-commit hooks on all files"
	@echo "  help - Show this help message"
