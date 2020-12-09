#/*******************************************************************************
# * Licensed Materials - Property of IBM
# * IBM Cloud Container Service, 5737-D43
# * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
# * US Government Users Restricted Rights - Use, duplication or
# * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
# ******************************************************************************/

GOPACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v /e2e | grep -v /volume-providers/softlayer/ | grep -v /samples) # With glide: GOPACKAGES=$(shell glide novendor)
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOLINTPACKAGES=$(shell go list ./... | grep -v /vendor/ | grep -v /e2e | grep -v /volume-providers/softlayer/ )
ARCH = $(shell uname -m)

export LINT_VERSION="1.27.0"

COLOR_YELLOW=\033[0;33m
COLOR_RESET=\033[0m

.PHONY: all
all: deps fmt vet test

.PHONY: deps
deps:
	# glide install
	go get github.com/pierrre/gotestcover
	@if ! which golangci-lint >/dev/null || [[ "$$(golangci-lint --version)" != *${LINT_VERSION}* ]]; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v${LINT_VERSION}; \
	fi

.PHONY: fmt
fmt: lint
	golangci-lint run --disable-all --enable=gofmt --timeout 600s --skip-dirs=e2e

.PHONY: dofmt
dofmt:
	golangci-lint run --disable-all --enable=gofmt --fix --skip-dirs=e2e

.PHONY: lint
lint:
	 golangci-lint run --skip-dirs=e2e

.PHONY: test
test:
ifeq ($(ARCH), ppc64le)
	# POWER
	$(GOPATH)/bin/gotestcover -v -coverprofile=cover.out ${GOPACKAGES} -timeout 90m
else
	# x86_64
	$(GOPATH)/bin/gotestcover -v -race -coverprofile=cover.out ${GOPACKAGES} -timeout 90m
endif

.PHONY: coverage
coverage:
	go tool cover -html=cover.out -o=cover.html

.PHONY: vet
vet:
	go vet ${GOPACKAGES}

.PHONY: build
build:
	go build -gcflags '-N -l' -o libSample samples/main.go samples/attach_detach.go samples/volume_operations.go

.PHONY: volume-lib-e2e-test
volume-lib-e2e-test:
	go test ./e2e/... -v -p 1 -ginkgo.progress -ginkgo.v -ginkgo.trace -timeout 240m  2>&1 | tee e2e_logs.txt

.PHONY: clean
clean:
	rm -rf libSample
