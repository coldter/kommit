.PHONY: build deps distclean
export LDFLAGS += -s
export LDFLAGS += -w

export CGO_ENABLED ?= 0
export GOCACHE ?= $(CURDIR)/.gocache

deps:
	go mod download

build: deps
	go build -o ./bin/ -ldflags='$(LDFLAGS)' ./...

install:
	go install -ldflags='$(LDFLAGS)' ./...

#lint:
#	@if ! command -v golangci-lint > /dev/null 2>&1; then \
#		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
#		sh -s -- -b "$$(go env GOPATH)/bin" v1.62.2 ; \
#	fi
#	golangci-lint run ./...

clean:
	go clean -x ./...

distclean:
	go clean -x -cache -testcache -modcache ./...
