# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOGET=$(GOCMD) get
BINARY_NAME=change-git-user

all: clean test vet build
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd
test:
	$(GOTEST) -v ./... -run ^Test -coverprofile cp.out
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd
	./$(BINARY_NAME)
vet:
	$(GOVET) ./cmd
deps:
	$(GOGET) github.com/markbates/goth
	$(GOGET) github.com/markbates/pop