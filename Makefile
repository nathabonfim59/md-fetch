.PHONY: build test clean install lint build-musl package-deb package-rpm packages build-all

BINARY_NAME=md-fetch
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go')
VERSION=$(shell git describe --tags --always || echo "dev")
COMMIT_HASH=$(shell git rev-parse --short HEAD || echo "unknown")

build: build-linux build-linux-musl

build-linux:
	GOOS=linux GOARCH=amd64 go build \
	-trimpath \
	-ldflags='-w -s' \
	-o $(BUILD_DIR)/$(BINARY_NAME)-$(COMMIT_HASH)-linux-amd64 ./cmd/fetch

build-linux-musl:
	CC=x86_64-linux-musl-gcc \
	CGO_ENABLED=1 \
	GOOS=linux GOARCH=amd64 \
	go build \
	-buildvcs=false \
	-trimpath \
	-ldflags='-w -s -linkmode external -extldflags "-static"' \
	-o $(BUILD_DIR)/$(BINARY_NAME)-$(COMMIT_HASH)-linux-musl-amd64 ./cmd/fetch

build-macos:
	GOOS=darwin GOARCH=amd64 go build \
	-trimpath \
	-ldflags='-w -s' \
	-o $(BUILD_DIR)/$(BINARY_NAME)-$(COMMIT_HASH)-darwin-amd64 ./cmd/fetch
	GOOS=darwin GOARCH=arm64 go build \
	-trimpath \
	-ldflags='-w -s' \
	-o $(BUILD_DIR)/$(BINARY_NAME)-$(COMMIT_HASH)-darwin-arm64 ./cmd/fetch

build-windows:
	GOOS=windows GOARCH=amd64 go build \
	-trimpath \
	-ldflags='-w -s' \
	-o $(BUILD_DIR)/$(BINARY_NAME)-$(COMMIT_HASH)-windows-amd64.exe ./cmd/fetch

build-all: build-linux build-linux-musl build-macos build-windows packages

package-deb: build-linux-musl
	ln -sf $(BINARY_NAME)-$(COMMIT_HASH)-linux-musl-amd64 $(BUILD_DIR)/$(BINARY_NAME)-package
	VERSION=$(VERSION) COMMIT_HASH=$(COMMIT_HASH) nfpm pkg --config build/nfpm.yml --target $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)+$(COMMIT_HASH)_amd64.deb --packager deb
	rm -f $(BUILD_DIR)/$(BINARY_NAME)-package

package-rpm: build-linux-musl
	ln -sf $(BINARY_NAME)-$(COMMIT_HASH)-linux-musl-amd64 $(BUILD_DIR)/$(BINARY_NAME)-package
	VERSION=$(VERSION) COMMIT_HASH=$(COMMIT_HASH) nfpm pkg --config build/nfpm.yml --target $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)+$(COMMIT_HASH).x86_64.rpm --packager rpm
	rm -f $(BUILD_DIR)/$(BINARY_NAME)-package

packages: package-deb package-rpm

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
	./$(BUILD_DIR)/$(BINARY_NAME)-$(COMMIT_HASH)-linux-amd64
