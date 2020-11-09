all: build

build:
	CGO_ENABLED=0 go build -tags netgo -ldflags '-extldflags "-static"'

package: build
	tar cjf notify-api.tar.bz2 notify-api

.PHONY: build
