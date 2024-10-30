# Copyright (c) Microsoft Corporation.
# Licensed under the Apache v2.0 License.
GOCMD=go
GOBUILD=$(GOCMD) build -v 
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)

BIN_DIR=bin
LD_FLAGS_WINDOWS_CSHARED=-extldflags=-Wl,--out-implib=MocCppWrapper.lib
CPP_WRAPPER_NAME=MocCppWrapper
CPP_WRAPPER_EXT=.dll
CPP_WRAPPER_OUT=$(BIN_DIR)/$(CPP_WRAPPER_NAME)$(CPP_WRAPPER_EXT)

# Private repo workaround
export GOPRIVATE=github.com/microsoft
# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

PKG := 

all: format  build unittest

clean:
	rm -rf ${OUT} ${OUTEXE} 

.PHONY: vendor
vendor:
	go mod tidy

build:
	GOARCH=amd64 go build -v ./...
	GOARCH=amd64 GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc $(GOBUILD) -buildmode=c-shared -o $(CPP_WRAPPER_OUT) -ldflags="$(LD_FLAGS_WINDOWS_CSHARED)" github.com/microsoft/moc-sdk-for-go/wrapper/cpp/

test:
	GOARCH=amd64 go test -v ./...

format:
	gofmt -s -w pkg/ services/ 

unittest:
	GOARCH=amd64 go test -v ./services/security/...
