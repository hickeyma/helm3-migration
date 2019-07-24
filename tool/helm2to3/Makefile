.PHONY: build
build: build-linux


build-linux: export GOOS=linux
build-linux: export GOARCH=amd64
build-linux: export CGO_ENABLED=0
build-linux: export GO111MODULE=on
build-linux:
	go build -v -o bin/helm2to3 cmd/helm2to3/*.go  # linux
