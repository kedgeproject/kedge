//
// Copyright (C) 2015-2017 Thomas LE ROUX <thomas@leroux.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package dockerparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShortParse(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "foo/bar")

	is.Equal("foo/bar:latest", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("latest", reference.Tag())
	is.Equal("docker.io", reference.Registry())
	is.Equal("docker.io/foo/bar", reference.Repository())
	is.Equal("docker.io/foo/bar:latest", reference.Remote())
}

func TestShortParseWithTag(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "foo/bar:1.1")

	is.Equal("foo/bar:1.1", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("1.1", reference.Tag())
	is.Equal("docker.io", reference.Registry())
	is.Equal("docker.io/foo/bar", reference.Repository())
	is.Equal("docker.io/foo/bar:1.1", reference.Remote())

}

func TestShortParseWithDigest(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb")

	is.Equal("foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb", reference.Tag())
	is.Equal("docker.io", reference.Registry())
	is.Equal("docker.io/foo/bar", reference.Repository())
	is.Equal("docker.io/foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb", reference.Remote())

}

func TestRegistry(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "localhost.localdomain/foo/bar")

	is.Equal("foo/bar:latest", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("latest", reference.Tag())
	is.Equal("localhost.localdomain", reference.Registry())
	is.Equal("localhost.localdomain/foo/bar", reference.Repository())
	is.Equal("localhost.localdomain/foo/bar:latest", reference.Remote())

}

func TestRegistryWithTag(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "localhost.localdomain/foo/bar:1.1")

	is.Equal("foo/bar:1.1", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("1.1", reference.Tag())
	is.Equal("localhost.localdomain", reference.Registry())
	is.Equal("localhost.localdomain/foo/bar", reference.Repository())
	is.Equal("localhost.localdomain/foo/bar:1.1", reference.Remote())

}

func TestRegistryWithDigest(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "localhost.localdomain/foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb")

	is.Equal("foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb", reference.Tag())
	is.Equal("localhost.localdomain", reference.Registry())
	is.Equal("localhost.localdomain/foo/bar", reference.Repository())
	is.Equal("localhost.localdomain/foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb", reference.Remote())

}

func TestRegistryWithPort(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "localhost.localdomain:5000/foo/bar")

	is.Equal("foo/bar:latest", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("latest", reference.Tag())
	is.Equal("localhost.localdomain:5000", reference.Registry())
	is.Equal("localhost.localdomain:5000/foo/bar", reference.Repository())
	is.Equal("localhost.localdomain:5000/foo/bar:latest", reference.Remote())

}

func TestRegistryWithPortAndTag(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "localhost.localdomain:5000/foo/bar:1.1")

	is.Equal("foo/bar:1.1", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("1.1", reference.Tag())
	is.Equal("localhost.localdomain:5000", reference.Registry())
	is.Equal("localhost.localdomain:5000/foo/bar", reference.Repository())
	is.Equal("localhost.localdomain:5000/foo/bar:1.1", reference.Remote())

}

func TestRegistryWithPortAndDigest(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "localhost.localdomain:5000/foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb")

	is.Equal("foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb", reference.Tag())
	is.Equal("localhost.localdomain:5000", reference.Registry())
	is.Equal("localhost.localdomain:5000/foo/bar", reference.Repository())
	is.Equal("localhost.localdomain:5000/foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb", reference.Remote())

}

func TestHttpRegistryClean(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "http://localhost.localdomain:5000/foo/bar:latest")

	is.Equal("foo/bar:latest", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("latest", reference.Tag())
	is.Equal("localhost.localdomain:5000", reference.Registry())
	is.Equal("localhost.localdomain:5000/foo/bar", reference.Repository())
	is.Equal("localhost.localdomain:5000/foo/bar:latest", reference.Remote())

}

func TestHttpsRegistryClean(t *testing.T) {

	is := require.New(t)

	reference := parse(is, "https://localhost.localdomain:5000/foo/bar:latest")

	is.Equal("foo/bar:latest", reference.Name())
	is.Equal("foo/bar", reference.ShortName())
	is.Equal("latest", reference.Tag())
	is.Equal("localhost.localdomain:5000", reference.Registry())
	is.Equal("localhost.localdomain:5000/foo/bar", reference.Repository())
	is.Equal("localhost.localdomain:5000/foo/bar:latest", reference.Remote())

}

func TestParseError(t *testing.T) {

	is := require.New(t)

	reference, err := Parse("sftp://user:passwd@example.com/foo/bar:latest")

	is.Error(err)
	is.Nil(reference)

}

func parse(is *require.Assertions, remote string) *Reference {

	reference, err := Parse(remote)

	is.NoError(err, "parse error was not expected")
	is.NotNil(reference)
	is.NotEmpty(reference)

	return reference
}
