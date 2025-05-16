OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
VERSION ?= $(shell git describe --tags --abbrev=0)
COMPILER_VERSION ?= $(shell go version | awk '{print $$3}')

PACK_NAME = gbase64
BIN_NAME = $(PACK_NAME)
PACK_PATH = bin/$(ARCH)/$(OS)

ifeq ($(OS), windows)
	BIN_NAME := $(BIN_NAME).exe
endif

.PHONY: bin
bin: vendor
	GOOS=$(OS) GOARCH=$(ARCH) go build \
		-ldflags "-s -w -X main.version=$(VERSION) -X main.compilerVersion=$(COMPILER_VERSION)" \
		-o $(PACK_PATH)/$(BIN_NAME) github.com/armourstill/$(PACK_NAME)

.PHONY: release
release: bin
	cd $(PACK_PATH) && zip $(PACK_NAME)-$(VERSION).$(OS).$(ARCH).zip $(BIN_NAME)

vendor:
	go mod tidy && go mod vendor
	find vendor -type d -exec chmod +w {} \;
