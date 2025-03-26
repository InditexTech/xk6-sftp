PROJECT_VERSION := 1.0.0

GOPATH := $(shell command go env GOPATH)

XK6_VERSION := v0.13.4
XK6_BINARY := "$(GOPATH)/bin/xk6"

GOLANGCI_VERSION := v1.64.5
GOLANGCI_BINARY := "$(GOPATH)/bin/golangci-lint"

.DEFAULT_GOAL := all

.PHONY: all
all: format lint compose-up test run compose-down

.PHONY: deps
deps:
	@if [ -z "$(XK6_BINARY)" ]; then \
		echo "Installing xk6..."; \
		go install go.k6.io/xk6/cmd/xk6@$(XK6_VERSION); \
	else \
		echo "xk6 is already installed."; \
	fi

	@if [ -z "$(GOLANGCI_BINARY)" ]; then \
			echo "Installing golangci-lint..."; \
			go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION); \
	else \
		echo "golangci-lint is already installed."; \
	fi

.PHONY: compose-up
compose-up:
	@echo "Starting sftp server..."
	@docker compose -f docker/docker-compose.yaml up -d

.PHONY: compose-down
compose-down:
	@echo "Destrying sftp server..."
	@docker compose -f docker/docker-compose.yaml down

.PHONY: build
build: deps
	@echo "Building k6  with STP extension..."
	@"$(XK6_BINARY)" build --with github.com/InditexTech/xk6-sftp=.

.PHONY: run
run: deps compose-up
	@echo "Running example..."
	@xk6 run --vus=1 --iterations=10 ./examples/main.js

.PHONY: verify
verify: compose-up deps format lint test compose-down

.PHONY: test
test:
	@echo "Running tests..."
	@go clean -testcache && go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: tidy
tidy:
	@echo "Running go mod tidy..."
	@go mod tidy

.PHONY: format
format:
	@echo "Running go fmt..."
	@go fmt ./...

.PHONY: lint
lint: deps
	@echo "Running golangci-lint..."
	@"$(GOLANGCI_BINARY)" run

.PHONY: reuse-deps
reuse-deps:
	@if [ -z "reuse" ]; then \
		echo "Installing reuse tool..."; \
		pip3 install --user reuse ;\
	else \
		echo "reuse is already installed."; \
	fi

.PHONY: add-copyright-headers
reuse-annotate: reuse-deps
	@echo "Adding copyright headers..."
	@reuse annotate --copyright "Industria de Diseño Textil S.A. INDITEX" --license "Apache-2.0" --year "$$(date +%Y)" --merge-copyrights *.go
	@reuse lint

.PHONY: get-version
get-version:
	@echo $(PROJECT_VERSION)
