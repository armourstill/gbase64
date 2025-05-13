OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
VERSION ?= $(shell git describe --tags --abbrev=0)
COMPILER_VERSION ?= $(shell go version | awk '{print $$3}')

.PHONY: bin
bin: vendor
	GOOS=$(OS) GOARCH=$(ARCH) go build \
		-ldflags "-X main.version=$(VERSION) -X main.compilerVersion=$(COMPILER_VERSION)" \
		-o bin/$(ARCH)/$(OS)/gbase64 github.com/armourstill/gbase64

.PHONY: bin.release
bin.release: vendor
	GOOS=$(OS) GOARCH=$(ARCH) go build \
		-ldflags "-w -s -X main.version=$(VERSION) -X main.compilerVersion=$(COMPILER_VERSION)" \
		-o bin/$(ARCH)/$(OS)/gbase64 github.com/armourstill/gbase64

vendor:
	go mod tidy && go mod vendor
	find vendor -type d -exec chmod +w {} \;
