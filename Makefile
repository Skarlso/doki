NAME=doki

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := bin
VERSION ?= v0.0.1

# List the GOOS and GOARCH to build
GO_LDFLAGS_STATIC="-s -w $(CTIMEVAR) -extldflags -static -X 'main.version=${VERSION}'"

.DEFAULT_GOAL := binaries

.PHONY: binaries
binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm darwin/amd64" \
		-ldflags=${GO_LDFLAGS_STATIC} \
		-output="$(BUILDDIR)/{{.OS}}/{{.Arch}}/$(NAME)" \
		-tags="netgo" \
		./

.PHONY: bootstrap
bootstrap:
	go get github.com/mitchellh/gox

.PHONY: lint
lint:
	golint ./...

.PHONY: run
run:
	go run main.go
