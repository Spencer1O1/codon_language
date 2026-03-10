GOCACHE ?= $(CURDIR)/.cache/go-build
GOMODCACHE ?= $(CURDIR)/.cache/go-mod

.PHONY: test build lint fmt vet clean tidy deps

test:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go test ./...

build:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build ./cmd/...

lint:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) golangci-lint run ./...

fmt:
	gofmt -w $(shell go list -f '{{.Dir}}' ./...)

vet:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go vet ./...

tidy:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go mod tidy

clean:
	rm -rf $(GOCACHE) $(GOMODCACHE)

deps:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go mod download
