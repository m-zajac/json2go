appname := json2go

sources := $(wildcard *.go)

build = GOOS=$(1) GOARCH=$(2) go build -o build/$(appname)$(3) ./cmd/json2go/*.go && chmod gu+x build/$(appname)$(3)
tar = cd build && tar -cvzf $(1)_$(2).tar.gz $(appname)$(3) && rm $(appname)$(3)
zip = cd build && zip $(1)_$(2).zip $(appname)$(3) && rm $(appname)$(3)

.PHONY: all windows darwin linux clean

all: windows darwin linux

clean:
	rm -rf build/

##### LINUX BUILDS #####
linux: build/linux_386.tar.gz build/linux_amd64.tar.gz

build/linux_386.tar.gz: $(sources)
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

##### DARWIN (MAC) BUILDS #####
darwin: build/darwin_amd64.tar.gz

build/darwin_amd64.tar.gz: $(sources)
	$(call build,darwin,amd64,)
	$(call tar,darwin,amd64)

##### WINDOWS BUILDS #####
windows: build/windows_386.zip build/windows_amd64.zip

build/windows_386.zip: $(sources)
	$(call build,windows,386,.exe)
	$(call zip,windows,386,.exe)

build/windows_amd64.zip: $(sources)
	$(call build,windows,amd64,.exe)
	$(call zip,windows,amd64,.exe)
