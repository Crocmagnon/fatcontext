.PHONY: lint
lint:
	pre-commit run --all-files

.PHONY: test
test:
	go test -race -v ./...
