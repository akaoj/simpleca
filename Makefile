.PHONY: help
.PHONY: clean compile release tests
.PHONY: _dev_image _reqs

.DEFAULT_GOAL = help

SHELL = /bin/bash
PWD = $(shell pwd)
CURRENT_USER_ID = $(shell id --user)
CURRENT_USER_NAME = $(shell id --user --name)

CONTAINER_IMAGE_NAME = simpleca:dev

define container_make
sudo docker run -it --rm -v ${PWD}:${PWD}:Z ${CONTAINER_IMAGE_NAME} "cd ${PWD} && make -f Makefile.container $1 $2"
endef

define HELP_CONTENT
Available `make` targets:
  clean    Clean the built binary, release package and Docker image.
  compile  Compile the code.
  help     Display this help.
  release  Generate a .tar.gz file for a release.
  tests    Run tests.
           Options: `IGNORE_VERSION=true` if you want to skip the version check (this can be annoying when developing
                    on a random branch).

Note that you don't need `go` or anything on your machine: all compilation, tests and packaging are done inside a
container. You can just `git clone` and `make release` out of the box (the only prerequisite is to be able to call
Docker with `sudo`).
endef
export HELP_CONTENT


VERSION = $(shell awk 'match($$0, /^const VERSION string = "(.*?)"$$/,a) {print a[1]}' src/main.go)
ARCH = $(shell uname -m)

IGNORE_VERSION = false

BINARY = simpleca-${VERSION}-${ARCH}


GO_SOURCES := $(shell find src/ -name '*.go')


_reqs:
	@command -v docker >/dev/null 2>&1 || { echo 'docker is needed'; exit 1; }

_dev_image: _reqs
	sudo docker build . --file dev.Dockerfile -t ${CONTAINER_IMAGE_NAME} --build-arg='USER_ID=${CURRENT_USER_ID}' --build-arg='USER_NAME=${CURRENT_USER_NAME}'


${BINARY}: ${GO_SOURCES}
	$(call container_make,compile,BINARY=${BINARY})


clean:
	$(RM) ${BINARY}
	$(RM) ${BINARY}.tar.gz
	-sudo docker rmi ${CONTAINER_IMAGE_NAME}


compile: _dev_image ${BINARY}


help:
	@echo -e "$${HELP_CONTENT}"


release: compile tests
	$(call container_make,release,BINARY=${BINARY})


tests: compile
	$(call container_make,tests,BINARY=${BINARY} VERSION=${VERSION} IGNORE_VERSION=${IGNORE_VERSION})
