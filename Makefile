 # Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    BINARY_NAME=dodas
    
    all: deps build

    build:
		$(GOBUILD) -o $(BINARY_NAME) -v

    test: deps build
		$(GOTEST) -v ./...
		echo "--- Valid template test ---"
		./dodas validate --template tests/tosca/valid_template.yml
		echo "--- Broken type test ---"
		./dodas validate --template tests/tosca/broken_template_type.yaml
		echo "--- Broken inputs in template test ---"
		./dodas validate --template tests/tosca/broken_template_node.yaml

    clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)

    run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

    deps:
		$(GOGET) github.com/spf13/cobra
		$(GOGET) github.com/spf13/viper
		$(GOGET) github.com/mitchellh/go-homedir
		$(GOGET) gopkg.in/yaml.v2
		$(GOGET) github.com/dciangot/toscalib
    
    docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/Cloud-PG/dodas-go-client golang:latest go build -o "$(BINARY_NAME)" -v