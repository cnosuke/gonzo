NAME     := gonzo
VERSION  := $(shell git describe --tags 2>/dev/null)
REVISION := $(shell git rev-parse --short HEAD 2>/dev/null)
SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

export GO111MODULE = on

bin/$(NAME): $(SRCS)
	go build -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(NAME) main.go

.PHONY: test deps build-for-docker build-docker push-docker

build-for-docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$(NAME)

deps:
	go get

test:
	go test -v ./...

build-docker: build-for-docker
	docker build -t cnosuke/gonzo:latest .

push-docker: build-docker
	docker push cnosuke/gonzo:latest
