GOCACHE ?= $(CURDIR)/.cache/go-build
GOMODCACHE ?= $(CURDIR)/.cache/go-mod
ROOT ?=

.PHONY: test build lint fmt vet clean tidy deps load validate validate-codon_language sync-assets sync-core-assets

test:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go test -count=1 ./cmd/... ./pkg/... ./internal/...

build:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go build ./cmd/...

load:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run ./cmd/codon load $(ROOT)

validate-codon_language:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run ./cmd/codon validate .codon
validate-codon:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run ./cmd/codon validate $(ROOT)
validate: test

emit:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run ./cmd/codon emit $(ROOT)

sync-core-assets:
	./scripts/sync_core_assets.sh $(ROOT)

sync-assets: sync-core-assets

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
