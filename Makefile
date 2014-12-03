SHELL = /bin/sh

GO_PACKAGES := $(shell go list ./... | sed 's_github.com/heroku/json2envdir_._')
GO_GIT_DESCRIBE_SYMBOL ?= github.com/heroku/json2envdir/config.version
GO_GIT_DESCRIBE := $(shell git describe --tags --always)
GO_BUILD_FLAGS := -ldflags "-X $(GO_GIT_DESCRIBE_SYMBOL) $(GO_GIT_DESCRIBE)"

all: build

build: godep
	godep go install $(GO_BUILD_FLAGS) ./...

imports:
	goimports -w $(GO_PACKAGES)

tidy: goimports
	test -z "$$(goimports -l -d $(GO_PACKAGES) | tee /dev/stderr)"

lint: golint
	test -z "$$(golint ./... | tee /dev/stderr)"

# dependencies
godep:
	go get github.com/tools/godep

goimports:
	go get code.google.com/p/go.tools/cmd/goimports

golint:
	go get github.com/golang/lint/golint


# build stuff

debversion   := $(shell cat config.version)
tempdir      := $(shell mktemp -d)
controldir   := $(tempdir)/DEBIAN
installpath  := $(tempdir)/usr/bin
buildpath    := .build
buildpackage := $(buildpath)/cache

define DEB_CONTROL
Package: json2envdir
Version: $(debversion)
Architecture: amd64
Maintainer: "Ricardo Chimal, Jr." <ricardo@heroku.com>
Section: heroku
Priority: optional
Description: JSON to envdir style directories
endef
export DEB_CONTROL

deb: bin/json2envdir
	echo "making deb"
	mkdir -p -m 0755 $(controldir)
	echo "$$DEB_CONTROL" > $(controldir)/control
	mkdir -p $(installpath)
	install bin/json2envdir $(installpath)/json2envdir
	fakeroot dpkg-deb -Z gzip --build $(tempdir) .
	rm -rf $(tempdir)

clean:
	rm -rf $(buildpath)
	rm -f json2envdir*.deb

bin/json2envdir:
	git clone git://github.com/heroku/heroku-buildpack-go.git $(buildpath)
	$(buildpath)/bin/compile . $(buildpackcache)

