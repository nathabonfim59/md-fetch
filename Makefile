.PHONY: build test clean install lint snapshot release release-major release-minor release-patch

BINARY_NAME=md-fetch
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go')

# Get the latest tag without the 'v' prefix
CURRENT_VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0")
# Split version into major, minor, and patch
MAJOR=$(shell echo $(CURRENT_VERSION) | cut -d. -f1)
MINOR=$(shell echo $(CURRENT_VERSION) | cut -d. -f2)
PATCH=$(shell echo $(CURRENT_VERSION) | cut -d. -f3)

build:
	goreleaser build --clean --single-target --snapshot

build-all:
	goreleaser build --clean --snapshot

snapshot:
	goreleaser release --clean --snapshot

release-major:
	@echo "Current version: v$(CURRENT_VERSION)"
	$(eval NEW_VERSION=$(shell echo $$(($(MAJOR)+1)).0.0))
	@echo "Creating major release v$(NEW_VERSION)"
	@read -p "Press enter to continue..." \
	&& git tag -a v$(NEW_VERSION) -m "Release v$(NEW_VERSION)" \
	&& git push origin v$(NEW_VERSION) \
	&& goreleaser release --clean

release-minor:
	@echo "Current version: v$(CURRENT_VERSION)"
	$(eval NEW_VERSION=$(shell echo $(MAJOR).$$(($(MINOR)+1)).0))
	@echo "Creating minor release v$(NEW_VERSION)"
	@read -p "Press enter to continue..." \
	&& git tag -a v$(NEW_VERSION) -m "Release v$(NEW_VERSION)" \
	&& git push origin v$(NEW_VERSION) \
	&& goreleaser release --clean

release-patch:
	@echo "Current version: v$(CURRENT_VERSION)"
	$(eval NEW_VERSION=$(shell echo $(MAJOR).$(MINOR).$$(($(PATCH)+1))))
	@echo "Creating patch release v$(NEW_VERSION)"
	@read -p "Press enter to continue..." \
	&& git tag -a v$(NEW_VERSION) -m "Release v$(NEW_VERSION)" \
	&& git push origin v$(NEW_VERSION) \
	&& goreleaser release --clean

release:
	@echo "Please use one of:"
	@echo "  make release-major  - for major version updates (X.0.0)"
	@echo "  make release-minor  - for minor version updates (x.Y.0)"
	@echo "  make release-patch  - for patch version updates (x.y.Z)"
	@echo "This ensures proper version management and tag creation."

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
