.PHONY: build test clean install lint build-musl

BINARY_NAME=fetch
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go')

build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(GO_FILES)
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/fetch

build-musl:
	CC=x86_64-linux-musl-gcc \
	CGO_ENABLED=1 \
	GOOS=linux GOARCH=amd64 \
	go build \
	-buildvcs=false \
	-trimpath \
	-ldflags='-w -s -linkmode external -extldflags "-static"' \
	-o $(BUILD_DIR)/$(BINARY_NAME)-musl ./cmd/fetch

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
