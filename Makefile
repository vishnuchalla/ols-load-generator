
.PHONY: build help images push


ARCH ?= amd64
BIN_NAME = ols-load-generator
BIN_DIR = bin
BIN_PATH = $(BIN_DIR)/$(ARCH)/$(BIN_NAME)
VERSION ?= $(shell git describe --tags --always)
CGO = 0

SOURCES := $(shell find . -type f -name "*.go")

# Containers
ENGINE ?= podman
REGISTRY = quay.io
ORG ?= vchalla
CONTAINER_NAME_ARCH = $(REGISTRY)/$(ORG)/ols-load-generator:$(ARCH)

all: lint build images push

help:
	@echo "Commands for $(BIN_PATH):"
	@echo
	@echo 'Usage:'
	@echo '    make clean                    Clean the compiled binaries'
	@echo '    [ARCH=arch] make build        Compile the project for arch, default amd64'
	@echo '    [ARCH=arch] make install      Installs ols-load-generator binary in the system, default amd64'
	@echo '    [ARCH=arch] make images       Build images for arch, default amd64'
	@echo '    [ARCH=arch] make push         Push images for arch, default amd64'
	@echo '    make help                     Show this message'

build: $(BIN_PATH)

$(BIN_PATH): $(SOURCES)
	@echo -e "\033[2mBuilding $(BIN_PATH)\033[0m"
	@echo "GOPATH=$(GOPATH)"
	GOARCH=$(ARCH) CGO_ENABLED=$(CGO) go build -v -o $(BIN_PATH) -ldflags "-X main.Version=$(VERSION)" ./cmd/ols-load-generator

lint:
	find . -name '*.go' -type f -exec go fmt {} \;

install:
	cp $(BIN_PATH) /usr/bin/$(BIN_NAME)

images:
	@echo -e "\033[2mBuilding container $(CONTAINER_NAME_ARCH)\033[0m"
	$(ENGINE) build --arch=$(ARCH) -f Containerfile $(BIN_DIR)/$(ARCH)/ -t $(CONTAINER_NAME_ARCH)

push:
	@echo -e "\033[2mPushing container $(CONTAINER_NAME_ARCH)\033[0m"
	$(ENGINE) push $(CONTAINER_NAME_ARCH)