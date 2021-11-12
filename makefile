# This makefile requires the go tools.
# This makefile also requires docker and docker-compose (https://www.docker.com/)
# See also subordinate makefile dependencies

# Parameters
BINARY_NAME=$(shell basename `pwd`)

# Commands
GOCMD=go
GOFMT=$(GOCMD) fmt
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

MAKE=make
RM=rm
DOCKERCOMPOSE=docker-compose

all: tools gen test install
tools:
	$(GOINSTALL) \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/golang/protobuf/protoc-gen-go \
gen:
	# the following makefiles have their own dependencies. See comments in each one.
	$(MAKE) -C frontend
	rm -rf assets/frontend && cp -R frontend/dist assets/frontend
	$(MAKE) -C assets
	$(MAKE) -C proto
	$(GOFMT) ./...
test: 
	$(GOTEST) -v ./...
install:
	$(GOINSTALL) -v
clean: 
	$(GOCLEAN) -v
	$(RM) -rf $(BINARY_LINUX) vendor output
	$(MAKE) -C frontend clean
	$(MAKE) -C assets clean
	$(MAKE) -C proto clean

linux: gen test
	@mkdir -p output/
	GOOS=linux GOARCH=amd64 go build -o "output/$(BINARY_NAME)" .
	@echo "see binary output/$(BINARY_NAME)"

docker: gen test
	$(DOCKERCOMPOSE) build
