.PHONY: build test clean install lint

BINARY_NAME=fetch
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go')

build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(GO_FILES)
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
	go fmt ./...
	go vet ./...

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)
