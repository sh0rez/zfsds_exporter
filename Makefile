VERSION := $(shell git describe --tags --dirty --always)

LDFLAGS := '-s -w -extldflags "-static" -X main.Version=${VERSION}'
static:
	CGO_ENABLED=0 go build -ldflags=${LDFLAGS} .
