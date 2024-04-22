APP_NAME=docker-bakery
VERSION=1.4.1

.DEFAULT_GOAL: all

.PHONY: all test build fmt install lint install-lint ci

all: install fmt build test lint

install:
	@echo "Installing dependencies"
	dep ensure -v

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

ci: build test lint

build:
	@echo "Building sources"
	go build -v ./...

release: fmt build test lint
	@echo $(VERSION)
	./release.sh $(VERSION)
