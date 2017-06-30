

GITCOMMIT := $(shell git rev-parse --short HEAD)
BUILD_FLAGS := -ldflags="-w -X github.com/kedgeproject/kedge/cmd.GITCOMMIT=$(GITCOMMIT)"
PKGS = $(shell glide novendor)

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
	go test $(PKGS)

# Run all tests
.PHONY: test
test: test-dep check-vendor validate test-unit

# Install all the required test-dependencies before executing tests (only valid when running `make test`)
.PHONY: test-dep
test-dep:
	go get github.com/mattn/goveralls
	go get github.com/modocache/gover
	go get github.com/Masterminds/glide
	go get github.com/sgotti/glide-vc
	go get github.com/golang/lint/golint
