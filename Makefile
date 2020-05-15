.PHONY: build build-docker serve update setup run-docker test coverage fmt doc

ifndef VERBOSE
MAKEFLAGS+=--no-print-directory
endif

ifeq ($(UNAME),Darwin)
ECHO=echo
else
ECHO=echo -e
endif

# Package
PACKAGE_NAME=go-microservice
PACKAGE_VERSION=0.0.1-alpha
BUILD=$(shell git rev-list --count HEAD)
ARCHITECTURE=amd64
# LDFLAGS=-ldflags '-v'
LDFLAGS=-ldflags '-w -s -v'

SRCS=./cmd/microservice/*.go

default: build

build:
	-@$(ECHO) "\n\033[0;35m%%% Building...\033[0m"
	-@$(ECHO) "Building..."
	CGO_ENABLED=0 go build $(LDFLAGS) -v -o ./dist/$(PACKAGE_NAME) $(SRCS)
	-@$(ECHO) "\n\033[1;32mDone!\033[0;32m\nDone!\033[0m\n"

build-docker:
	docker build -t $(PACKAGE_NAME) .

serve:
	go run ./cmd/microservice/main.go serve

update:
	go run ./cmd/microservice/main.go update

setup:
	go run ./cmd/microservice/main.go setup

run-docker:
	docker run -d -p 3035:3035 $(PACKAGE_NAME)

test:
	-@$(ECHO) "\n\033[0;35m%%% Running tests\033[0m"
	go test -v ./...

coverage:
	-@$(ECHO) "\n\033[0;35m%%% Running test coverage\033[0m"
	go test -cover ./...

doc:
  godoc -http=:6060 -index

fmt:
	go fmt ./...