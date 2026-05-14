# iwmo Go library — development tasks.
#
# Default target runs the fast unit tests. Integration tests need a live iWMO
# endpoint and stay off by default; run `make test-integration` with the
# IWMO_* env vars set to exercise them.

GO          ?= go
GOLANGCI    ?= golangci-lint
FUZZTIME    ?= 30s

.PHONY: all build test test-race test-integration test-coverage coverage-html lint fmt vet fuzz clean

all: test

build:
	$(GO) build ./...

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

test-integration:
	$(GO) test -tags integration -v ./...

test-coverage:
	$(GO) test -coverprofile=coverage.out ./... && $(GO) tool cover -func=coverage.out | tail -1

coverage-html:
	$(GO) test -coverprofile=coverage.out ./... && $(GO) tool cover -html=coverage.out -o coverage.html

lint:
	$(GOLANGCI) run ./...

fmt:
	gofmt -s -w .

vet:
	$(GO) vet ./...

fuzz:
	$(GO) test -run=^$$ -fuzz=FuzzDecode -fuzztime=$(FUZZTIME) ./...

clean:
	rm -f coverage.out coverage.html
