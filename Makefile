# Main targets for a Go app project
#
# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

OS = $(shell uname | tr A-Z a-z)
export PATH := $(abspath bin/):${PATH}

# Build variables
BUILD_DIR ?= build
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DATE_FMT = +%FT%T%z
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif
LDFLAGS += -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}
export CGO_ENABLED ?= 0
ifeq (${VERBOSE}, 1)
ifeq ($(filter -v,${GOARGS}),)
	GOARGS += -v
endif
TEST_FORMAT = short-verbose
endif

# Dependency versions
GOTESTSUM_VERSION ?= 0.4.2
GOLANGCI_VERSION ?= 1.27.0

GOLANG_VERSION ?= 1.14

.PHONY: clear
clear: ${CLEAR_TARGETS} ## Clear the working area and the project
	rm -rf bin/

.PHONY: clean
clean: ${CLEAN_TARGETS} ## Clean builds
	rm -rf ${BUILD_DIR}/

.PHONY: run-%
run-%: build-%
	${BUILD_DIR}/$*

.PHONY: run
run: $(patsubst cmd/%,run-%,$(wildcard cmd/*)) ## Build and execute all applications

.PHONY: goversion
goversion:
ifneq (${IGNORE_GOLANG_VERSION}, 1)
	@printf "${GOLANG_VERSION}\n$$(go version | awk '{sub(/^go/, "", $$3);print $$3}')" | sort -t '.' -k 1,1 -k 2,2 -k 3,3 -g | head -1 | grep -q -E "^${GOLANG_VERSION}$$" || (printf "Required Go version is ${GOLANG_VERSION}\nInstalled: `go version`" && exit 1)
endif

.PHONY: build-deps
build-deps: ${BUILD_DEP_TARGETS}
	@:

.PHONY: pre-build
pre-build: ${PRE_BUILD_TARGETS}
	@:

.PHONY: post-build
post-build: ${POST_BUILD_TARGETS}
	@:

.PHONY: build-%
build-%: build-deps pre-build
build-%: goversion
ifeq (${VERBOSE}, 1)
	go env
endif

	@mkdir -p ${BUILD_DIR}
	go build ${GOARGS} -trimpath -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/$* ./cmd/$*

	@${MAKE} post-build

.PHONY: build
build: build-deps pre-build
build: goversion ## Build binaries
ifeq (${VERBOSE}, 1)
	go env
endif

	@mkdir -p ${BUILD_DIR}
	go build ${GOARGS} -trimpath -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/ ./cmd/...

	@${MAKE} post-build

.PHONY: build-release-deps
build-release-deps: build-deps
build-release-deps: ${BUILD_RELEASE_DEP_TARGETS}
	@:

.PHONY: pre-build-release
pre-build-release: ${PRE_BUILD_RELEASE_TARGETS}
	@:

.PHONY: post-build-release
post-build-release: ${POST_BUILD_RELEASE_TARGETS}
	@:

.PHONY: build-release
build-release: build-release-deps pre-build-release
build-release: ## Build binaries without debug information
	@${MAKE} LDFLAGS="-w ${LDFLAGS}" GOARGS="${GOARGS} -trimpath" BUILD_DIR="${BUILD_DIR}/release" build

	@${MAKE} post-build-release

.PHONY: build-debug-deps
build-debug-deps: build-deps
build-debug-deps: ${BUILD_DEBUG_DEP_TARGETS}
	@:

.PHONY: pre-build-debug
pre-build-debug: ${PRE_BUILD_DEBUG_TARGETS}
	@:

.PHONY: post-build-debug
post-build-debug: ${POST_BUILD_DEBUG_TARGETS}
	@:

.PHONY: build-debug
build-debug: build-debug-deps pre-build-debug
build-debug: ## Build binaries with remote debugging capabilities
	@${MAKE} GOARGS="${GOARGS} -gcflags \"all=-N -l\"" BUILD_DIR="${BUILD_DIR}/debug" build

	@${MAKE} post-build-debug

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum
bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_${OS}_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}

TEST_PKGS ?= ./...
.PHONY: test
test: TEST_FORMAT ?= short
test: SHELL = /bin/bash
test: export CGO_ENABLED=1
test: bin/gotestsum ## Run tests
	@mkdir -p ${BUILD_DIR}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/coverage.xml --format ${TEST_FORMAT} -- -race -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Variable outputting/exporting rules
var-%: ; @echo $($*)
varexport-%: ; @echo $*=$($*)
