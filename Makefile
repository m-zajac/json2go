appname := json2go

goroot := $(shell go env GOROOT)

build = GOOS=$(1) GOARCH=$(2) go build -o build/$(appname)$(3) ./cmd/json2go/*.go && chmod gu+x build/$(appname)$(3)
tar = cd build && tar -cvzf $(1)_$(2).tar.gz $(appname)$(3) && rm $(appname)$(3)
zip = cd build && zip $(1)_$(2).zip $(appname)$(3) && rm $(appname)$(3)
last_version = $(shell git describe --tags --abbrev=0)

.PHONY: all windows darwin linux web clean test lint lint-more depl-pages

all: windows darwin linux web

clean:
	rm -rf build/

test:
	go test -race ./...

lint: $(shell go env GOPATH)/bin/golangci-lint
	$(shell go env GOPATH)/bin/golangci-lint run ./...


##### DEPS #####

$(shell go env GOPATH)/bin/golint:
	GOFLAGS="-mod=readonly" go get golang.org/x/lint/golint

$(shell go env GOPATH)/bin/golangci-lint:
	GOFLAGS="-mod=readonly" go get github.com/golangci/golangci-lint/cmd/golangci-lint


##### BUILDS #####

linux: build/linux_386.tar.gz build/linux_amd64.tar.gz
darwin: build/darwin_amd64.tar.gz
windows: build/windows_386.zip build/windows_amd64.zip
web:
	mkdir -p build/web
	cp "$(goroot)/misc/wasm/wasm_exec.js" build/web
	GOOS=js GOARCH=wasm go build -o build/web/json2go.wasm ./cmd/json2go-ws/*.go

build/linux_386.tar.gz:
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz:
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

build/darwin_amd64.tar.gz:
	$(call build,darwin,amd64,)
	$(call tar,darwin,amd64)

build/windows_386.zip:
	$(call build,windows,386,.exe)
	$(call zip,windows,386,.exe)

build/windows_amd64.zip:
	$(call build,windows,amd64,.exe)
	$(call zip,windows,amd64,.exe)


##### DEPLOYMENTS #####

depl-pages: web
	mkdir -p build/gh-pages
	cp -r deployments/gh-pages/* build/gh-pages
	cp -r build/web/* build/gh-pages
	sed -i 's!_VERSION_!'"$(last_version)"'!g' build/gh-pages/index.html
