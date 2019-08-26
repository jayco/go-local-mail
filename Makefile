.PHONY: build buildlinux

buildlinux:
	GO111MODULE=on GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/linux/local-mail ./cmd/main.go

buildDarwin:
	GO111MODULE=on GOARCH=amd64 GOOS=darwin go build -ldflags="-w -s" -o bin/mac/local-mail ./cmd/main.go

build: buildlinux buildDarwin
