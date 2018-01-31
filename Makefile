GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean -x
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

SOURCE_FOLDERS?=$$(go list ./... | grep -v /vendor/)
TEST_PATTERN?=.
TEST_OPTIONS?=-race -covermode=atomic -coverprofile=coverage.txt
OUTPUT_NAME!=basename $$(pwd)

GIT_BRANCH!=git rev-parse --abbrev-ref HEAD
GIT_COMMIT!=git rev-parse --short HEAD
BUILD!=date +%FT%T%z
VERSION!=cat dist-version

LDFLAGS=-ldflags "-w -s \
				-X github.com/duckpuppy/${OUTPUT_NAME}/cmd.Version=${VERSION} \
				-X github.com/duckpuppy/${OUTPUT_NAME}/cmd.Build=${BUILD} \
				-X github.com/duckpuppy/${OUTPUT_NAME}/cmd.Commit=${GIT_COMMIT} \
				-X github.com/duckpuppy/${OUTPUT_NAME}/cmd.Branch=${GIT_BRANCH}"

setup: ## Install all the build and lint dependencies
	$(GOGET) -u github.com/alecthomas/gometalinter
	$(GOGET) -u github.com/golang/dep/...
	$(GOGET) -u github.com/pierrre/gotestcover
	$(GOGET) -u golang.org/x/tools/cmd/cover
	$(GOGET) -u github.com/mitchellh/gox
	dep ensure
	gometalinter --install --update

clean: ## Clean up build artifacts
	$(GOCLEAN)
	rm -Rf build
	rm -Rf dist
	rm -f coverage.txt

test: ## Run all the tests
	gotestcover $(TEST_OPTIONS) $(SOURCE_FOLDERS) -run $(TEST_PATTERN) -timeout=1m

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run all the linters
	gometalinter --vendor --disable-all \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gofmt \
		--enable=goimports \
		--enable=dupl \
		--enable=misspell \
		--enable=errcheck \
		--enable=vet \
		--enable=vetshadow \
		--deadline=10m \
		./...

ci: lint test ## Run all the tests and code checks

build: ## Build a beta version
	$(GOBUILD) ${LDFLAGS} -race -o ./build/$(OUTPUT_NAME) -v

install: ## Install to $GOPATH/src
	go install ${LDFLAGS} ./...

dist: ## Build all the distribution files
	gox ${LDFLAGS} -output="./dist/${OUTPUT_NAME}_{{.OS}}_{{.Arch}}" -verbose

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
