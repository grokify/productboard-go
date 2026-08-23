.PHONY: all build test test-coverage lint generate update-spec clean tidy fmt vet check help

all: build

# Build all packages
build:
	go build ./...

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run ./...

# Regenerate internal/api from the currently-vendored openapi/entities.yaml.
# Does NOT re-fetch the spec — see update-spec.
generate:
	./generate.sh

# Deliberately re-fetch the latest ProductBoard OpenAPI spec. Never run as
# part of generate/build — spec changes should be reviewed before codegen.
# Despite the .yaml URL, the response body is actually JSON (confirmed
# 2026-08-12) — saved as .json to match its real content type. Update
# openapi/SOURCE.md's fetch date after running this.
update-spec:
	curl -sL https://developer.productboard.com/openapi/entities.yaml -o openapi/entities.json
	@echo "Vendored spec updated. Review the diff, update openapi/SOURCE.md, then run 'make generate'."

# Clean build artifacts
clean:
	rm -rf internal/api/*.go
	rm -f coverage.out coverage.html

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	gofmt -s -w .

# Run go vet
vet:
	go vet ./...

# Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

help:
	@echo "Available targets:"
	@echo "  build          - Build all packages"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run golangci-lint"
	@echo "  generate       - Regenerate internal/api from the vendored spec"
	@echo "  update-spec    - Re-fetch the ProductBoard OpenAPI spec (review before generate)"
	@echo "  check          - Run fmt, vet, lint, test"
