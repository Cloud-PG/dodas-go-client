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
		./dodas --config ${HOME}/.dodas_go_client.yaml validate --template ${HOME}/git/TOSCA_BARI/_htcondor_.yml
    clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
    run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
    deps:
		$(GOGET) github.com/spf13/cobra
		$(GOGET) github.com/spf13/viper
		$(GOGET) github.com/dciangot/toscalib
    
    docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/cloudpg/dodas-go-client golang:latest go build -o "$(BINARY_NAME)" -v