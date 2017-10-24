# Copyright 2017 The Kedge Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GITCOMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
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
	./scripts/run_e2e.sh -p $(PARALLEL) -t $(TIMEOUT)
else
ifdef PARALLEL
	./scripts/run_e2e.sh -p $(PARALLEL)
else
ifdef TIMEOUT
	./scripts/run_e2e.sh -t $(TIMEOUT)
else
	./scripts/run_e2e.sh
endif
endif
endif

# Run all tests
.PHONY: test
test: test-dep validate test-unit test-unit-cover

# Tests that are run on travs-ci
.PHONY: travis-tests
travis-tests: test-dep validate test-unit-cover test-jsonschema-generation

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

# Update vendoring
# Vendoring is a bit messy right now
.PHONY: vendor-update
vendor-update:
	# Handles packages defined in glide.yaml
	glide update -v
	# Vendors OpenShift and its dependencies
	./scripts/vendor-openshift.sh

# Test if the changed types.go is valid and JSONSchema can be generated out of it
.PHONY: test-jsonschema-generation
test-jsonschema-generation:
	docker run -v `pwd`/pkg/spec/types.go:/data/types.go:ro,Z surajd/kedgeschema

