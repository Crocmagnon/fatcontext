.PHONY: lint
lint:
	pre-commit run --all-files golangci-lint-full
