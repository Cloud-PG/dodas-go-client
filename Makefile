 # Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    BINARY_NAME=dodas_go_client
    
    all: test build
    build: 
		$(GOBUILD) -o $(BINARY_NAME) -v
    test: 
		$(GOTEST) -v ./...
    clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
    run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
    deps:
		$(GOGET) github.com/spf13/cobra
		$(GOGET) github.com/spf13/viper
    $(GOGET) github.com/owulveryck/toscalib
    
    
    # Cross compilation
    build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
    docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/cloudpg/dodas_go_client golang:latest go build -o "$(BINARY_NAME)" -v