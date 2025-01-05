
    .PHONY: build test clean

    BINARY_NAME=fetch
    BUILD_DIR=bin

    build:
        mkdir -p $(BUILD_DIR)
        go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/fetch

    test:
        go test ./...

    clean:
        rm -rf $(BUILD_DIR)
        go clean

    install:
        go install ./cmd/fetch

    lint:
        golangci-lint run

    .DEFAULT_GOAL := build
