VERSION?=`git describe --tags`
DOCBIN?=mkdocs
BUILD_DATE := `date +%Y-%m-%d\ %H:%M`
VERSIONFILE := version.go

GOCMD=go
GOBUILD=$(GOCMD) build -x -ldflags "-w -v"
GOBUILD_DBG=$(GOCMD) build -x
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
	$(GOBUILD) -o $(BINARY_NAME)

build-debug:
	$(GOBUILD_DBG) -o $(BINARY_NAME) -v

doc:
	cp README.md docs/README.md
	BUILD_DOC=true ./$(BINARY_NAME)

publish-doc:
	$(DOCBIN) gh-deploy

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
	docker run --rm -it -v ${PWD}:/go -w /go/ golang:1.12.1 go build -o "$(BINARY_NAME)" -v

docker-img-build:
	docker build . -t dodas

windows-build:
	env GOOS=windows $(GOBUILD) -o $(BINARY_NAME).exe -v

macos-build:
	env GOOS=darwin $(GOBUILD) -o $(BINARY_NAME)_osx -v

gensrc:
	rm -f $(VERSIONFILE)
	@echo "package main" > $(VERSIONFILE)
	@echo "const (" >> $(VERSIONFILE)
	@echo "  VERSION = \"$(VERSION)\"" >> $(VERSIONFILE)
	@echo "  BUILD_DATE = \"$(BUILD_DATE)\"" >> $(VERSIONFILE)
	@echo ")" >> $(VERSIONFILE)

build-release: tidy gensrc build doc publish-doc test windows-build macos-build
	zip dodas.zip dodas
	zip dodas.exe.zip dodas.exe
	zip dodas_osx.zip dodas_osx
