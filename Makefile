NAME=avahi-discovery
LDFLAGS="-w -s"
SRC=main.go inventory.go

build: fmt
	@mkdir -p bin/
	GO111MODULE=on GOSUMDB=off GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${NAME}-linux-amd64 -ldflags ${LDFLAGS} ${SRC}
	GO111MODULE=on GOSUMDB=off GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/${NAME}-darwin-amd64 -ldflags ${LDFLAGS} ${SRC}

fmt:
	go fmt ./...

.PHONY: fmt