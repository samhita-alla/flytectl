export REPOSITORY=flytectl
include boilerplate/lyft/golang_test_targets/Makefile

GIT_VERSION := $(shell git describe --always --tags)
GIT_HASH := $(shell git rev-parse --short HEAD)
TIMESTAMP := $(shell date '+%Y-%m-%d')
PACKAGE ?=github.com/flyteorg/flytestdlib

LD_FLAGS="-s -w -X $(PACKAGE)/version.Version=$(GIT_VERSION) -X $(PACKAGE)/version.Build=$(GIT_HASH) -X $(PACKAGE)/version.BuildTime=$(TIMESTAMP)"



define PIP_COMPILE
pip-compile $(1) --upgrade --verbose
endef

generate:
	go test github.com/flyteorg/flytectl/cmd --update

compile:
	go build -o bin/flytectl -ldflags=$(LD_FLAGS) main.go

.PHONY: update_boilerplate
update_boilerplate:
	@boilerplate/update.sh

.PHONY: install-piptools
install-piptools:
	pip install -U pip-tools

.PHONY: doc-requirements.txt
doc-requirements.txt: doc-requirements.in install-piptools
	$(call PIP_COMPILE,doc-requirements.in)
