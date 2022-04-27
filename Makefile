NOW=`date '+%Y.%m.%d %H:%M:%S'`
OS=`uname -i -v`
AFTER_COMMIT=`git rev-parse HEAD`
VERSION=0.3.0

build:
	go build -ldflags "-X 'main.BuildVersion=$(VERSION)' -X 'main.BuildTime=$(NOW)' -X 'main.BuildOSUname=$(OS)' -X 'main.BuildCommit=$(AFTER_COMMIT)'" -o bin/phpsmith ./cmd/phpsmith

.PHONY: build
