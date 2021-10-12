CC_DIR := cc
DQL_Y := $(CC_DIR)/dql.y
DQL_GO := $(CC_DIR)/dql.go
DQL_OUTPUT := $(CC_DIR)/dql.output

generate: $(DQL_GO) go-generate

regenerate: clean generate

go-regenerate: clean-go-generate go-generate

.PHONY: go-generate
go-generate:
	go generate ./...

$(DQL_GO): $(DQL_Y)
	goyacc -o $(DQL_GO) -v $(DQL_OUTPUT) $(DQL_Y)

.PHONY: clean
clean: clean-dql clean-go-generate clean-dist

.PHONY: clean-dql
clean-dql:
	rm -f $(DQL_GO) $(DQL_OUTPUT)

.PHONY: clean-go-generate
clean-go-generate:
	find $(ROOT) -name "*_generated.go" -type f | xargs rm -f

.PHONY: test
test:
	go test ./...

.PHONY: prepare
prepare:
	go install github.com/berquerant/marker@latest
	go install github.com/berquerant/mkvisitor@latest
	go install golang.org/x/tools/cmd/stringer@latest

.PHONY: clean-dist
clean-dist:
	rm -rf dist

dist/debugger:
	go build -o $@ ./cmd/debugger

REPO := github.com/berquerant/dql
META_PKG := $(REPO)
COMMIT := $(shell git rev-parse HEAD)
TAG := $(shell git describe --tags)
GO_VERSION := $(shell go version)
LDFLAGS := "-X '$(META_PKG).Version=$(TAG)' -X '$(META_PKG).GoVersion=$(GO_VERSION)' -X '$(META_PKG).Commit=$(COMMIT)'"

dist/dql:
	go build -o $@ -ldflags $(LDFLAGS) ./cmd/dql
