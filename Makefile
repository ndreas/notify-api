all: build

build:
	CGO_ENABLED=0 go build -tags netgo -ldflags '-extldflags "-static"'

.PHONY: build
