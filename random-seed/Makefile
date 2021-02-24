#!/usr/bin/make -f

export GO111MODULE = on

install:
	@echo "installing..."
	@go build -o ${GOPATH}/bin/random-sp ./cmd

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs goimports -w -local github.com/irisnet/service-providers-go/random-seed
