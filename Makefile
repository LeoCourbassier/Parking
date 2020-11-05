GO ?= go
TEST_RUN ?= br.com.mlabs/models br.com.mlabs/usecases br.com.mlabs/api
GOBUILD ?= $(GO) build
RED=\033[0;31m
GREEN=\033[0;32m
NC=\033[0m
CODEZERO=$$?

.PHONY: all dep build clean test lint cover migration

all: clean build

lint:
	@echo "Linting" 
	@go vet -v ./... 
	@golint -set_exit_status ./...

test:
	@echo "Testing..." 
	$(GO) test $(TEST_RUN) -v -coverprofile coverage.out
ifeq ($$?, $(CODEZERO))
	@echo "${GREEN}Sucess"
endif

race: dep
	@$(GO) test -race -short ./...

msan: dep 
	@$(GO) test -msan -short ./...

dep:
	@echo "Getting the dependencies" 
	@$(GO) get -v -d ./...

build: dep
	@echo "Building server" 
	@mkdir -p cmd && $(GOBUILD) -o ./cmd ./...

clean:
	@echo "Cleaning the last build" 
	@rm -rf ./cmd

cover:
	@$(GO) tool cover -func=coverage.out