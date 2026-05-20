VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0-dev")
LDFLAGS := -ldflags="-X github.com/hantsaniala/dj/cmd.Version=$(VERSION)"

.PHONY: build release install clean

build:
	go build $(LDFLAGS) -o dj .

release:
	goreleaser release --clean

install:
	go install $(LDFLAGS) github.com/hantsaniala/dj@latest

clean:
	rm -f dj
