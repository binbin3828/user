# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=UserServer
BINARY_UNIX=$(BINARY_NAME)_unix

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
# windows
# $(GOTEST) -v ..\test\...
# linux
	$(GOTEST) -v ../test/...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Cross compilation, when on window to build linux bin file
build-linux:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

