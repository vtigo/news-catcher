.PHONY: all build clean test fmt vet run-dev run-api run-fetch run-cron run-tui

# Binary names
BINARY_DIR=bin
DEV_BINARY=$(BINARY_DIR)/dev
API_BINARY=$(BINARY_DIR)/api
FETCH_BINARY=$(BINARY_DIR)/fetch
CRON_BINARY=$(BINARY_DIR)/cron
TUI_BINARY=$(BINARY_DIR)/tui

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod

all: clean build

build: build-dev build-api build-fetch build-cron build-tui

build-dev:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(DEV_BINARY) ./cmd/dev

build-api:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(API_BINARY) ./cmd/api

build-fetch:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(FETCH_BINARY) ./cmd/fetch

build-cron:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(CRON_BINARY) ./cmd/cron

build-tui:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(TUI_BINARY) ./cmd/tui

run-dev: build-dev
	$(DEV_BINARY)

run-api: build-api
	$(API_BINARY)

run-fetch: build-fetch
	$(FETCH_BINARY)

run-cron: build-cron
	$(CRON_BINARY)

run-tui: build-tui
	$(TUI_BINARY)

test:
	$(GOTEST) -v ./...

fmt:
	$(GOFMT) ./...

vet:
	$(GOVET) ./...

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

deps:
	$(GOMOD) download
	$(GOMOD) tidy

help:
	@echo "Available targets:"
	@echo "  all         - Clean and build all binaries"
	@echo "  build       - Build all binaries"
	@echo "  build-dev   - Build dev binary"
	@echo "  build-api   - Build api binary"
	@echo "  build-fetch - Build fetch binary"
	@echo "  build-cron  - Build cron binary"
	@echo "  build-tui   - Build tui binary"
	@echo "  run-dev     - Build and run dev binary"
	@echo "  run-api     - Build and run api binary"
	@echo "  run-fetch   - Build and run fetch binary"
	@echo "  run-cron    - Build and run cron binary"
	@echo "  run-tui     - Build and run tui binary"
	@echo "  test        - Run tests"
	@echo "  fmt         - Format code"
	@echo "  vet         - Run go vet"
	@echo "  clean       - Remove binaries and clean"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  help        - Show this help message"
