VERSION=0.0.2
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
GO111MODULE=on

all: check-memcached-uptime

.PHONY: check-memcached-uptime

check-memcached-uptime: check-memcached-uptime.go
	go build $(LDFLAGS) -o check-memcached-uptime

linux: check-memcached-uptime.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o check-memcached-uptime

clean:
	rm -rf check-memcached-uptime

check:
	go test ./...

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin master
