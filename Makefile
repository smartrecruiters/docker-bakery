APP_NAME=docker-bakery
VERSION=1.0.5

.DEFAULT_GOAL: all

.PHONY: all test build fmt install

all: install fmt build test

install:
	@echo "Installing dependencies"
	go get -v ./...

fmt:
	@echo "Formating source code"
	goimports -l -w .

test:
	@echo "Running tests"
	go test -v ./... && echo "TESTS PASSED"

build: fmt test
	@echo "Building sources"
	go build -v ./...

release: build
	@echo $(VERSION)
	./release.sh $(VERSION)
