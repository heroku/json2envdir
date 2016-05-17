SHELL = /bin/sh

GO_LINKER_SYMBOL := main.version
GO_BUILD_ENV := GOOS=linux GOARCH=amd64

GO_PACKAGES := $(shell go list ./... | sed 's_github.com/heroku/json2envdir_._')

all: build

build:
	go install $(GO_BUILD_FLAGS) ./...

imports:
	goimports -w $(GO_PACKAGES)

tidy: goimports
	test -z "$$(goimports -l -d $(GO_PACKAGES) | tee /dev/stderr)"

lint: golint
	test -z "$$(golint ./... | tee /dev/stderr)"

goimports:
	go get code.google.com/p/go.tools/cmd/goimports

golint:
	go get github.com/golang/lint/golint


deb: tmp ldflags ver
	echo "making deb"
	$(eval DEB_ROOT := "${TMP}/DEBIAN")
	${GO_BUILD_ENV} go build -v -o ${TMP}/usr/bin/json2envdir ${LDFLAGS} ./cmd/json2envdir
	mkdir -p ${DEB_ROOT}
	cat misc/DEBIAN.control | sed s/{{VERSION}}/${VERSION}/ > ${DEB_ROOT}/control
	dpkg-deb -Zgzip -b ${TMP} json2envdir_${VERSION}_amd64.deb
	rm -rf ${TMP}

tmp:
	$(eval TMP := $(shell mktemp -d -t json2envdir.XXXXX))

ldflags: glv
	$(eval LDFLAGS := -ldflags "-X ${GO_LINKER_SYMBOL}=${GO_LINKER_VALUE}")

glv:
	$(eval GO_LINKER_VALUE := $(shell git describe --tags --always))

ver: glv
	$(eval VERSION := $(shell echo ${GO_LINKER_VALUE} | sed s/^v//))

clean:
	rm -f json2envdir*.deb
