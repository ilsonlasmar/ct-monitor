.PHONY: build

GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 sam build
