APP_NAME=codon
CMD_PATH=./cmd/codon
COMMAND?=
GENOME?=./examples/issue-tracker

.PHONY: dev run build install fmt vet test  lint clean

dev: fmt vet lint genome validate

run:
	go run $(CMD_PATH) $(COMMAND) $(GENOME)

genome:
	go run $(CMD_PATH) genome $(GENOME)

validate:
	go run $(CMD_PATH) validate $(GENOME)

build:
	go build -o $(APP_NAME) $(CMD_PATH)

install:
	go install $(CMD_PATH)

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -f $(APP_NAME)