.PHONY: lint
lint:
	pre-commit run --all-files

.PHONY: test
test:
	go mod download
	go test -race -v ./...
