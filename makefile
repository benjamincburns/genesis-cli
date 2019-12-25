GO111MODULE=on

DIRECTORIES=$(sort $(dir $(wildcard pkg/*/*/)))

MOCKS=$(foreach x, $(DIRECTORIES), mocks/$(x))

OUTPUT_DIR=./bin

.PHONY: build test test_race lint vet get mocks clean-mocks manual-mocks
.ONESHELL:

all: genesis

genesis: | prep get
	go build -o bin/genesis ./cmd/genesis

prep:
	@mkdir $(OUTPUT_DIR) 2>> /dev/null | true 

install:
	go install cmd/genesis/main.go

test:
	go test ./...

lint:
	golint ./...

vet:
	go vet ./...

get:
	go get ./...

clean-mocks:
	rm -rf mocks

mocks: $(MOCKS)
	
$(MOCKS): mocks/% : %
	mockery -output=$@ -dir=$^ -all
