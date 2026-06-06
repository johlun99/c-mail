.PHONY: lint test build dev fmt install-hooks

# Go packages excluding anything under frontend/ (e.g. stray .go files in node_modules)
GO_PKGS := $(shell go list ./... | grep -v /frontend/)

# Run all linters (Go + frontend)
lint:
	golangci-lint run
	cd frontend && npm run lint && npm run format:check && npm run typecheck

# Run all tests (Go + frontend)
test:
	go test $(GO_PKGS)
	cd frontend && npm run test

# Build a distributable binary
build:
	wails build

# Run the app in development with hot reload
dev:
	wails dev

# Auto-format Go + frontend
fmt:
	gofmt -w .
	cd frontend && npm run format

# Install git hooks (pre-commit lint + test)
install-hooks:
	lefthook install
