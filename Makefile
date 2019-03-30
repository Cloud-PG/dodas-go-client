 # Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    BINARY_NAME=dodas
    REPO=github.com/Cloud-PG/dodas-go-client

		export GO111MODULE=on

    all: build test

    build:
		$(GOBUILD) -o $(BINARY_NAME) -v

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
    
    docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/$(REPO) golang:latest go build -o "$(BINARY_NAME)" -v