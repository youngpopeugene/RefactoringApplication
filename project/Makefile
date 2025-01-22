LINTER=golangci-lint
GOSEC=gosec
TESTS=go test ./...

install-tools:
	@echo "Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest

format:
	@gofmt -w .
	@goimports -w .
	@echo "Code formatted successfully."

lint:
	@echo "Running linter..."
	@$(LINTER) run

test:
	@echo "Running tests..."
	@$(TESTS)

gosec:
	@echo "Running gosec..."
	@$(GOSEC) ./...

all: format lint test gosec
	@echo "All checks passed successfully."

.PHONY: install-tools format lint test gosec all
