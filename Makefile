VERSION=0.0.1

all: check-memcached-uptime

.PHONY: check-memcached-uptime

gom:
	go get -u github.com/mattn/gom

bundle:
	gom install

check-memcached-uptime: check-memcached-uptime.go
	gom build -o check-memcached-uptime

linux: check-memcached-uptime.go
	GOOS=linux GOARCH=amd64 gom build -o check-memcached-uptime

fmt:
	go fmt ./...

dist:
	git archive --format tgz HEAD -o check-memcached-uptime-$(VERSION).tar.gz --prefix check-memcached-uptime-$(VERSION)/

clean:
	rm -rf check-memcached-uptime check-memcached-uptime-*.tar.gz

