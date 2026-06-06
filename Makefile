.PHONY: lint test build dev fmt install-hooks

# Go packages excluding anything under frontend/ (e.g. stray .go files in node_modules)
GO_PKGS := $(shell go list ./... | grep -v /frontend/)

# Linux needs the webkit2gtk-4.1 build tag (Arch ships 4.1, not 4.0). macOS uses
# native WebKit and needs no tag.
WAILS_TAGS := $(if $(filter Linux,$(shell uname -s)),-tags webkit2_41,)

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
	wails build $(WAILS_TAGS)

# Run the app in development with hot reload
dev:
	wails dev $(WAILS_TAGS)

# Auto-format Go + frontend
fmt:
	gofmt -w .
	cd frontend && npm run format

# Install git hooks (pre-commit lint + test)
install-hooks:
	lefthook install
