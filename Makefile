PKG_PREFIX := github.com/imorph/blank

APP_NAME := blank

BUILDINFO_TAG ?= $(shell echo $$(git describe --long --all | tr '/' '-')$$( \
	      git diff-index --quiet HEAD -- || echo '-dirty-'$$(git diff-index -u HEAD | sha1sum | grep -oP '^.{8}')))

PKG_TAG ?= $(shell git tag -l --points-at HEAD)
ifeq ($(PKG_TAG),)
PKG_TAG := $(BUILDINFO_TAG)
endif


GO_BUILDINFO = -X '$(PKG_PREFIX)/lib/buildinfo.Version=$(APP_NAME)-$(shell date -u +'%Y%m%d-%H%M%S')-$(BUILDINFO_TAG)'

all: clean deps check_all build

build:
	GO111MODULE=on go build .

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