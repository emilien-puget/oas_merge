all: gofumpt import lint

init:
	go install mvdan.cc/gofumpt@v0.6.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2
	go install github.com/daixiang0/gci@v0.12.3

lint:
	golangci-lint run ./...

gofumpt:
	gofumpt -l -w .

import:
	gci write --skip-generated -s standard -s default -s "prefix(github.com/emilien-puget/oas_merge)" -s blank -s dot -s alias .
