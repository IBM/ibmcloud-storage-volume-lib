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

.PHONY: all
all: deps fmt vet test

.PHONY: deps
deps:
	glide install
	go get github.com/pierrre/gotestcover
	go get golang.org/x/lint/golint

.PHONY: fmt
fmt: lint
	gofmt -l ${GOFILES}
	@if [ -n "$$(gofmt -l ${GOFILES})" ]; then echo 'Above Files needs gofmt fixes. Please run gofmt -l -w on your code.' && exit 1; fi

.PHONY: lint
lint:
	$(GOPATH)/bin/golint --set_exit_status ${GOLINTPACKAGES}

.PHONY: makefmt
makefmt:
	gofmt -l -w ${GOFILES}

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
	go build -gcflags '-N -l' -o libSample samples/main.go samples/attach_detach.go

.PHONY: volume-lib-e2e-test
volume-lib-e2e-test:
	go test ./e2e/... -v -p 1 -ginkgo.progress -ginkgo.v -ginkgo.trace -ginkgo.failFast true -timeout 90m

.PHONY: clean
clean:
	rm -rf libSample
