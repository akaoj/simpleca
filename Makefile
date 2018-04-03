.PHONY: build
.PHONY: reqs_build  reqs_tests

SHELL := /bin/bash

TESTS_DIR := tests
SIMPLECA := simpleca


GO_SOURCES := $(shell find . -name '*.go')


reqs_build:
	@command -v go >/dev/null 2>&1 || { echo 'go binary is needed'; exit 1; }

reqs_tests: reqs_build
	@command -v openssl >/dev/null 2>&1 || { echo 'openssl binary is needed'; exit 1; }


build: reqs_build ${SIMPLECA}

${SIMPLECA}: ${GO_SOURCES}
	go build -o ${SIMPLECA}


include Makefile.tests
