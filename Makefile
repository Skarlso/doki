NAME=doki

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := bin
VERSION ?= v0.0.1

# List the GOOS and GOARCH to build
GO_LDFLAGS_STATIC="-s -w $(CTIMEVAR) -extldflags -static -X 'main.version=${VERSION}'"

.DEFAULT_GOAL := binaries

##@ Build

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
lint: ## Run the linter
	golangci-lint run --exclude-use-default=false --timeout=5m0s

.PHONY: run
run:
	go run main.go

##@ Tests

unit: ## Run the unit tests
	ginkgo -r ./pkg

##@ Docs

docs: mdtoc ## Update the Readme
	mdtoc -inplace README.md

mdtoc: ## Download mdtoc binary if necessary
	GO111MODULE=off go get sigs.k8s.io/mdtoc || true
