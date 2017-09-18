

GITCOMMIT := $(shell git rev-parse --short HEAD)
BUILD_FLAGS := -ldflags="-w -X github.com/kedgeproject/kedge/cmd.GITCOMMIT=$(GITCOMMIT)"
PKGS = $(shell glide novendor)
UNITPKGS = $(shell glide novendor | grep -v tests)

default: bin

.PHONY: all
all: bin

.PHONY: bin
bin:
	go build ${BUILD_FLAGS} -o kedge main.go

.PHONY: install
install:
	go install ${BUILD_FLAGS}

# kompile kedge for multiple platforms
.PHONY: cross
cross:
	go get github.com/mitchellh/gox
	gox -osarch="darwin/amd64 linux/amd64 linux/arm windows/amd64" -output="bin/kedge-{{.OS}}-{{.Arch}}" $(BUILD_FLAGS)

.PHONY: clean
clean:
	rm -f kedge
	rm -r -f bundles

# run all validation tests
.PHONY: validate
validate: gofmt vet lint

.PHONY: vet
vet:
	go vet $(PKGS)

# golint errors are only recommendations
.PHONY: lint
lint:
	golint $(PKGS)

.PHONY: gofmt
gofmt:
	./scripts/check-gofmt.sh

# Checks if there are nested vendor dirs inside Kompose vendor and if vendor was cleaned by glide-vc
.PHONY: check-vendor
check-vendor:
	./scripts/check-vendor.sh

.PHONY: test-unit
test-unit:
	go test $(UNITPKGS)

.PHONY: test-e2e
test-e2e:

ifneq ($(and $(PARALLEL),$(TIMEOUT)),)
	go test -parallel=$(PARALLEL) -timeout=$(TIMEOUT) -v github.com/kedgeproject/kedge/tests/e2e
else
ifdef PARALLEL
	go test -parallel=$(PARALLEL) -v github.com/kedgeproject/kedge/tests/e2e
else
ifdef TIMEOUT
	go test -timeout=$(TIMEOUT) -v github.com/kedgeproject/kedge/tests/e2e
else
	go test -v github.com/kedgeproject/kedge/tests/e2e
endif
endif
endif

# Run all tests
.PHONY: test
test: test-dep check-vendor validate test-unit

# Tests that are run on travs-ci
.PHONY: travis-tests
travis-tests: test-dep check-vendor validate test-unit-cover

# Install all the required test-dependencies before executing tests (only valid when running `make test`)
.PHONY: test-dep
test-dep:
	go get github.com/mattn/goveralls
	go get github.com/modocache/gover
	go get github.com/Masterminds/glide
	go get github.com/sgotti/glide-vc
	go get github.com/golang/lint/golint

# Run unit tests and collect coverage
.PHONY: test-unit-cover
test-unit-cover:
	# First install packages that are dependencies of the test.
	go test -i -race -cover $(PKGS)
	# go test doesn't support colleting coverage across multiple packages,
	# generate go test commands using go list and run go test for every package separately
	go list -f '"go test -race -cover -v -coverprofile={{.Dir}}/.coverprofile {{.ImportPath}}"' github.com/kedgeproject/kedge/...  | grep -v "vendor" | grep -v "e2e" | xargs -L 1 -P4 sh -c
