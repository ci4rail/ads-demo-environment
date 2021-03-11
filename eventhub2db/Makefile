BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/ads-demo-environment/eventhub2db/cmd.version=$(VERSION)"

build:
	GOOS=linux GOARCH=amd64 go build $(GO_LDFLAGS) -o ${BIN_DIR}/eventhub2db main.go

test:
	go test ./...
clean:
	rm -f ${BIN_DIR}/eventhub2db

.PHONY: clean
