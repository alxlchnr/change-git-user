# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOGET=$(GOCMD) get
GOGENERATE=$(GOCMD) generate
BINARY_NAME=change-git-user

all: deps clean generate test vet build
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
generate:
	$(GOGENERATE) ./...
deps:
	$(GOGET) github.com/golang/mock/gomock
	$(GOGET) github.com/golang/mock/mockgen
