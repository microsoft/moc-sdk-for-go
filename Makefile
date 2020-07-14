# Copyright (c) Microsoft Corporation.
# Licensed under the Apache v2.0 License.
GOCMD=go
GOBUILD=$(GOCMD) build -v 
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)
# Private repo workaround
export GOPRIVATE = github.com/microsoft
# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

PKG := 

all: format  build

clean:
	rm -rf ${OUT} ${OUTEXE} 

.PHONY: vendor
vendor:
	go mod tidy

build:
	GOARCH=amd64 go build -v ./...

test:
	GOARCH=amd64 go test -v ./...

format:
	gofmt -s -w pkg/ services/ 
