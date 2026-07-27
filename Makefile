VERSION=0.0.13
GITCOMMIT?=$(shell git describe --dirty --always)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.commit=${GITCOMMIT}"

all: check-lastlog

check-lastlog: main.go
	go build $(LDFLAGS) -o check-lastlog

linux: main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o check-lastlog

check:
	go test -v ./...
