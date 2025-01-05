.PHONY: build test clean run install lint

BINARY_NAME := md-fetch
BUILD_DIR := bin
GO := go
GOLANGCI_LINT := golangci-lint

build:
	mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/fetch

run: build
	$(BUILD_DIR)/$(BINARY_NAME)

test:
	$(GO) test ./...

clean:
	rm -rf $(BUILD_DIR)
	$(GO) clean

install:
	$(GO) install ./cmd/fetch

lint:
	$(GOLANGCI_LINT) run

.DEFAULT_GOAL := build
