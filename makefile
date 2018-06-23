# This makefile requires the go tools and dep, which can be installed on a Mac with Homebrew:
# brew install go dep
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

DEP=dep
MAKE=make
RM=rm
DOCKERCOMPOSE=docker-compose

all: generate vendor test install
generate:
	# the following makefiles have their own dependencies. See comments in each one.
	$(MAKE) -C frontend
	rm -rf bindata/frontend && cp -R frontend/dist bindata/frontend
	$(MAKE) -C bindata
	$(MAKE) -C proto
	$(MAKE) -C mocks
	$(GOFMT) ./...
vendor:
	$(DEP) ensure
test: 
	$(GOTEST) -v ./...
install:
	$(GOINSTALL) -v
clean: 
	$(GOCLEAN) -v
	$(RM) -rf $(BINARY_LINUX) vendor output
	$(MAKE) -C frontend clean
	$(MAKE) -C bindata clean
	$(MAKE) -C proto clean
	$(MAKE) -C mocks clean

linux: generate vendor test
	@mkdir -p output/
	GOOS=linux GOARCH=amd64 go build -o "output/$(BINARY_NAME)" .
	@echo "see binary output/$(BINARY_NAME)"

docker: generate vendor test
	$(DOCKERCOMPOSE) build