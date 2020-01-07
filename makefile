GO111MODULE=on
CGO_ENABLED=0


OUTPUT_DIR=./bin
INSTALL_LOC=/usr/local/bin

SHORT_HASH=$(shell git log --pretty=format:'%h' -n 1)
DATE=$(shell date +"%d.%m.%y")

LDFLAGS=-X main.buildTime=$(DATE) -X main.commitHash=$(SHORT_HASH)

BUILD_FLAGS=-tags netgo -ldflags '$(LDFLAGS) -extldflags "-static"'
LINUX_FLAGS=$(BUILD_FLAGS)
MAC_FLAGS=$(BUILD_FLAGS)
WINDOWS_FLAGS=$(BUILD_FLAGS)

PKG=main

.PHONY: build test lint vet get linux darwin windows multiplatform install clean freebsd version
.ONESHELL:

all: genesis

genesis: | prep get
	@go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/genesis ./cmd/genesis

clean:
	rm -rf $(OUTPUT_DIR)/genesis
	rm -rf $(OUTPUT_DIR)/linux
	rm -rf $(OUTPUT_DIR)/darwin
	rm -rf $(OUTPUT_DIR)/windows
	rm -rf $(OUTPUT_DIR)/freebsd

prep:
	@mkdir $(OUTPUT_DIR) 2>> /dev/null || true

linux:
	@mkdir -p $(OUTPUT_DIR)/linux 2>> /dev/null || true
	GOARCH=386 GOOS=linux go build $(LINUX_FLAGS) -o $(OUTPUT_DIR)/linux/386/genesis ./cmd/genesis
	GOARCH=amd64 GOOS=linux go build $(LINUX_FLAGS) -o $(OUTPUT_DIR)/linux/amd64/genesis ./cmd/genesis
	GOARCH=arm GOOS=linux go build $(LINUX_FLAGS) -o $(OUTPUT_DIR)/linux/arm/genesis ./cmd/genesis
	GOARCH=arm64 GOOS=linux go build $(LINUX_FLAGS) -o $(OUTPUT_DIR)/linux/arm64/genesis ./cmd/genesis
	GOARCH=ppc64 GOOS=linux go build $(LINUX_FLAGS) -o $(OUTPUT_DIR)/linux/ppc64/genesis ./cmd/genesis

freebsd:
	@mkdir -p $(OUTPUT_DIR)/freebsd 2>> /dev/null || true
	GOARCH=386 GOOS=freebsd go build $(LINUX_FLAGS) -o $(OUTPUT_DIR)/freebsd/386/genesis ./cmd/genesis
	GOARCH=amd64 GOOS=freebsd go build $(LINUX_FLAGS) -o $(OUTPUT_DIR)/freebsd/amd64/genesis ./cmd/genesis
	GOARCH=arm GOOS=freebsd go build $(LINUX_FLAGS) -o $(OUTPUT_DIR)/freebsd/arm/genesis ./cmd/genesis

darwin:
	@mkdir -p $(OUTPUT_DIR)/darwin 2>> /dev/null || true
	GOARCH=amd64 GOOS=darwin go build $(MAC_FLAGS) -o $(OUTPUT_DIR)/darwin/amd64/genesis ./cmd/genesis

windows:
	@mkdir -p $(OUTPUT_DIR)/windows 2>> /dev/null || true
	GOARCH=amd64 GOOS=windows go build $(WINDOWS_FLAGS) -o $(OUTPUT_DIR)/windows/amd64/genesis.exe ./cmd/genesis 

multiplatform: linux darwin windows freebsd

install: | genesis
	@cd cmd/genesis && \
	go install -ldflags="$(LDFLAGS)" . &&\
	cd -

test:
	go test ./...

lint:
	golint ./...

vet:
	go vet ./...

get:
	@go get ./...

version:
	@echo -n $(SHORT_HASH)