
# Copyright 2016 The Kubernetes Authors.
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


GITCOMMIT := $(shell git rev-parse --short HEAD)
BUILD_FLAGS := -ldflags="-w -X github.com/surajssd/kapp/cmd.GITCOMMIT=$(GITCOMMIT)"
PKGS = $(shell glide novendor)

default: bin

.PHONY: all
all: bin

.PHONY: bin
bin:
	go build ${BUILD_FLAGS} -o kapp main.go

.PHONY: install
install:
	go install ${BUILD_FLAGS}

# kompile kapp for multiple platforms
.PHONY: cross
cross:
	gox -osarch="darwin/amd64 linux/amd64 linux/arm windows/amd64" -output="bin/kapp-{{.OS}}-{{.Arch}}" $(BUILD_FLAGS)

.PHONY: clean
clean:
	rm -f kapp
	rm -r -f bundles

# run all validation tests
.PHONY: validate
validate: gofmt vet lint

.PHONY: vet
vet:
	go vet $(PKGS)

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

# Run all tests
.PHONY: test
test: test-dep check-vendor validate install

# Install all the required test-dependencies before executing tests (only valid when running `make test`)
.PHONY: test-dep
test-dep:
	go get github.com/mattn/goveralls
	go get github.com/modocache/gover
	go get github.com/Masterminds/glide
	go get github.com/sgotti/glide-vc
	go get github.com/golang/lint/golint
