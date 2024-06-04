#! /usr/bin/make

# 确定目标平台
TARGET_OS ?= linux
TARGET_ARCH ?= amd64

# 二进制名字
ifndef BINARY_NAME
	BINARY_NAME=bypctl-$(TARGET_OS)-$(TARGET_ARCH)
endif

# 编译参数
LDFLAGS ?= -ldflags "-extldflags -static"

goreleaser:
	goreleaser release "--clean" "--snapshot"

build:
	CGO_ENABLED=0 GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build $(LDFLAGS) -o $(BINARY_NAME)

init:
	rm -rf go.mod go.sum vendor
	go mod init bypctl
	go mod edit -replace=go.opentelemetry.io/otel/sdk@v1.14.0=go.opentelemetry.io/otel/sdk@v1.24.0
	go mod tidy
	go mod vendor

run:
	@echo "bypctl are running"
	go run main.go

clean:
	go clean -i .
	rm -rf ${BINARY_NAME} go.mod go.sum docs vendor

help:
	@echo "make build: compile packages and dependencies"
	@echo "make init: download all packages from go.mod"
	@echo "make run: to run ```go run main.go server```"
	@echo "make clean: remove object files and cached files"
	@echo "make goreleaser: use goreleaser to build"
