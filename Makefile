PKG_PREFIX := github.com/imorph/blank

APP_NAME := blank

BUILDINFO_SHA ?= $(shell git rev-parse HEAD | grep -oP '^.{8}')
BUILDINFO_TAG ?= $(shell git describe --abbrev=0)
GO_BUILDINFO = -X '$(PKG_PREFIX)/lib/buildinfo.Version=$(APP_NAME)-$(BUILDINFO_TAG)-$(BUILDINFO_SHA)'

all: clean deps check_all build

build:
	GO111MODULE=on go build -ldflags "$(GO_BUILDINFO)" .

clean:
	rm -rf ./blank

check_all: fmt vet lint errcheck golangci-lint

deps:
	go get ./...

fmt:
	GO111MODULE=on gofmt -l -w -s .
	GO111MODULE=on gofmt -l -w -s .

lint: install-golint
	golint .

install-golint:
	which golint || GO111MODULE=off go get -u github.com/golang/lint/golint

vet:
	GO111MODULE=on go vet .

errcheck: install-errcheck
	errcheck .

install-errcheck:
	which errcheck || GO111MODULE=off go get -u github.com/kisielk/errcheck

golangci-lint: install-golangci-lint
	golangci-lint run -D errcheck

install-golangci-lint:
	which golangci-lint || GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint