OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
VERSION ?= $(shell git describe --tags --abbrev=0)
COMPILER_VERSION ?= $(shell go version | awk '{print $$3}')

BIN_PATH = bin/$(ARCH)/$(OS)
BIN_NAME = gbase64

.PHONY: bin
bin: vendor
	GOOS=$(OS) GOARCH=$(ARCH) go build \
		-ldflags "-s -w -X main.version=$(VERSION) -X main.compilerVersion=$(COMPILER_VERSION)" \
		-o $(BIN_PATH)/$(BIN_NAME) github.com/armourstill/$(BIN_NAME)

.PHONY: release
release: bin
	zip $(BIN_PATH)/$(BIN_NAME)-$(VERSION).$(OS).$(ARCH).zip $(BIN_PATH)/$(BIN_NAME)

vendor:
	go mod tidy && go mod vendor
	find vendor -type d -exec chmod +w {} \;
