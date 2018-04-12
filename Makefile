.PHONY: build
.PHONY: reqs_build  reqs_tests

SHELL := /bin/bash

TESTS_DIR := tests
SIMPLECA := simpleca


GO_SOURCES := $(shell find . -name '*.go')

VERSION := $(shell awk 'match($$0, /^const VERSION string = "(.*?)"$$/,a) {print a[1]}' main.go)
ARCH := $(shell uname -m)


reqs_build:
	@command -v go >/dev/null 2>&1 || { echo 'go binary is needed'; exit 1; }

reqs_tests: reqs_build
	@command -v openssl >/dev/null 2>&1 || { echo 'openssl binary is needed'; exit 1; }

reqs_release: reqs_build
	@command -v tar >/dev/null 2>&1 || { echo 'tar command is needed'; exit 1; }


build: reqs_build ${SIMPLECA}

${SIMPLECA}: ${GO_SOURCES}
	go build -o ${SIMPLECA}


release: reqs_release tests build
	cp ${SIMPLECA} ${SIMPLECA}-${VERSION}-${ARCH}
	tar --create --gzip --file ${SIMPLECA}-${VERSION}-${ARCH}.tar.gz ${SIMPLECA}-${VERSION}-${ARCH}


include Makefile.tests
