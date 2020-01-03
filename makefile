GO111MODULE=on
CGO_ENABLED=0


OUTPUT_DIR=./bin

SHORT_HASH=$(shell git log --pretty=format:'%h' -n 1)
DATE=$(shell date +"%d.%m.%y")

LDFLAGS=-X main.buildTime=$(DATE) -X main.commitHash=$(SHORT_HASH)

BUILD_FLAGS=-tags netgo -ldflags '$(LDFLAGS) -extldflags "-static"'
LINUX_FLAGS=$(BUILD_FLAGS)
MAC_FLAGS=$(BUILD_FLAGS)
WINDOWS_FLAGS=$(BUILD_FLAGS)

.PHONY: build test lint vet get linux mac windows multiplatform install clean
.ONESHELL:

all: genesis

genesis: | prep get
	@go build -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/genesis ./cmd/genesis

clean:
	rm -rf $(OUTPUT_DIR)/genesis
	rm -rf $(OUTPUT_DIR)/linux
	rm -rf $(OUTPUT_DIR)/mac
	rm -rf $(OUTPUT_DIR)/windows

prep:
	@mkdir $(OUTPUT_DIR) 2>> /dev/null | true

linux:
	@mkdir -p $(OUTPUT_DIR)/linux 2>> /dev/null | true
	GOOS=linux
	@go build $(LINUX_FLAGS) -o $(OUTPUT_DIR)/linux/genesis ./cmd/genesis

mac:
	@mkdir -p $(OUTPUT_DIR)/mac 2>> /dev/null | true
	GOARCH=amd64 GOOS=darwin go build $(MAC_FLAGS) -o $(OUTPUT_DIR)/mac/genesis ./cmd/genesis

windows:
	@mkdir -p $(OUTPUT_DIR)/windows 2>> /dev/null | true
	GOARCH=amd64 GOOS=windows go build $(WINDOWS_FLAGS) -o $(OUTPUT_DIR)/windows/genesis.exe ./cmd/genesis 

multiplatform: linux mac windows

install: | genesis
	go install cmd/genesis/main.go

test:
	go test ./...

lint:
	golint ./...

vet:
	go vet ./...

get:
	@go get ./...
