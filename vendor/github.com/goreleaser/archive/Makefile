SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

setup:
	go get -u golang.org/x/tools/cmd/cover
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go get -u -t  ./...

test:
	go test $(TEST_OPTIONS) -v -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m

cover: test
	go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m

fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint:
	./bin/golangci-lint run --enable-all ./...

ci: lint test

.DEFAULT_GOAL := build
