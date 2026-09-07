VERSION=0.0.20
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"

all: check-lastlog

check-lastlog: main.go
	go build $(LDFLAGS) -o check-lastlog

linux: main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o check-lastlog

check:
	go test -v ./...

lint:
	golangci-lint run ./...
