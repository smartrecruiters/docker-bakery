APP_NAME=docker-bakery
VERSION=1.0.6

.DEFAULT_GOAL: all

.PHONY: all test build fmt install lint install-lint

all: install fmt build test lint

install:
	@echo "Installing dependencies"
	go get -v ./...

fmt:
	@echo "Formating source code"
	goimports -l -w .

install-lint:
	@echo "Installing golinter"
	go get -u golang.org/x/lint/golint

lint:
	@echo "Executing golint"
	golint bakery/...

test:
	@echo "Running tests"
	go test -v ./... && echo "TESTS PASSED"

build: fmt test
	@echo "Building sources"
	go build -v ./...

release: build
	@echo $(VERSION)
	./release.sh $(VERSION)
