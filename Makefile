.PHONY: build test clean install lint snapshot release

BINARY_NAME=md-fetch
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go')

build:
	goreleaser build --clean --single-target --snapshot

build-all:
	goreleaser build --clean --snapshot

snapshot:
	goreleaser release --clean --snapshot

release:
	goreleaser release --clean

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -rf dist/
	go clean

install:
	go install ./cmd/fetch

lint:
	go fmt ./...
	go vet ./...

run: build
	./dist/md-fetch_*/md-fetch
