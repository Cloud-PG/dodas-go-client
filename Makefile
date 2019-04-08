GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=dodas
REPO=github.com/Cloud-PG/dodas-go-client

export GO111MODULE=on
# Force 64 bit architecture
export GOARCH=amd64

all: build test

build:
	$(GOBUILD) -o $(BINARY_NAME) -ldflags=-w -v

build-debug:
	$(GOBUILD) -o $(BINARY_NAME) -v

build-doc: build
	BUILD_DOC=true ./$(BINARY_NAME)

test: build
	$(GOTEST) -v ./...
	./$(BINARY_NAME) validate --template tests/tosca/valid_template.yml
	./$(BINARY_NAME) validate --template tests/tosca/broken_template_type.yaml
	./$(BINARY_NAME) validate --template tests/tosca/broken_template_node.yaml

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

install:
	$(GOCMD) install $(REPO)

tidy:
	$(GOCMD) mod tidy

docker-bin-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/$(REPO) golang:1.12.1 go build -o "$(BINARY_NAME)" -v

docker-img-build:
	docker build . -t dodas

windows-build:
	env GOOS=windows $(GOBUILD) -o $(BINARY_NAME).exe -v

macos-build:
	env GOOS=darwin $(GOBUILD) -o $(BINARY_NAME)_osx -v

build-release: tidy build test windows-build macos-build docker-img-build
	zip dodas.zip dodas
	zip dodas.exe.zip dodas.exe
	zip dodas_osx.zip dodas_osx
