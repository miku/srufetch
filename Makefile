SHELL = /bin/bash
TARGETS = srufetch
PKGNAME = srufetch

VERSION := $(shell ([ -f VERSION ] && cat VERSION) || (git rev-parse --short HEAD))
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

.PHONY: all assets bench clean clean-docs cloc deb deps imports lint members names rpm test vet

all: deps $(TARGETS)

deps:
	go get -v ./...

$(TARGETS): %: cmd/%/main.go
	go build -ldflags="-linkmode=external -X main.Version=$(VERSION) -X main.BuildTime=$(BUILDTIME)" -o $@ $<

clean:
	rm -f $(TARGETS)
	rm -f $(PKGNAME)_*deb
	rm -f $(PKGNAME)-*rpm
	rm -rf ./packaging/deb/$(PKGNAME)/usr

imports:
	go get golang.org/x/tools/cmd/goimports
	goimports -w .

# Packaging related.
deb: all
	mkdir -p packaging/deb/$(PKGNAME)/usr/local/bin
	cp $(TARGETS) packaging/deb/$(PKGNAME)/usr/local/bin
	cd packaging/deb && fakeroot dpkg-deb --build $(PKGNAME) .
	mv packaging/deb/$(PKGNAME)_*.deb .

rpm: all
	mkdir -p $(HOME)/rpmbuild/{BUILD,SOURCES,SPECS,RPMS}
	cp ./packaging/rpm/$(PKGNAME).spec $(HOME)/rpmbuild/SPECS
	cp $(TARGETS) $(HOME)/rpmbuild/BUILD
	./packaging/rpm/buildrpm.sh $(PKGNAME)
	cp $(HOME)/rpmbuild/RPMS/x86_64/$(PKGNAME)*.rpm .

