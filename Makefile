NOW=`date '+%Y.%m.%d %H:%M:%S'`
OS=`uname`
AFTER_COMMIT=`git rev-parse HEAD`
VERSION=0.3.0

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X 'main.BuildVersion=$(VERSION)' -X 'main.BuildTime=$(NOW)' -X 'main.BuildOSUname=$(OS)' -X 'main.BuildCommit=$(AFTER_COMMIT)'" -o bin/phpsmith ./cmd/phpsmith

.PHONY: build lint test ci-lint install-linter

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --config=.golangci.yml

test:
	@echo "Running tests..."
	@go test ./... -cover -short -count=1 -race

ci-lint: install-linter lint

install-linter:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
