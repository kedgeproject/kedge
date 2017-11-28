/*
Copyright 2017 The Kedge Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package build

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	dockerlib "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
)

func BuildPushDockerImage(dockerfile, image, context string, push bool) error {

	context, err := filepath.Abs(context)
	if err != nil {
		return err
	}
	// Connect to the Docker client
	client, err := dockerlib.NewClientFromEnv()
	if err != nil {
		return err
	}

	build := Build{Client: *client}
	if err = build.BuildImage(dockerfile, image, context); err != nil {
		return err
	}

	if push {
		if err = PushImage(image); err != nil {
			return err
		}
	}

	return nil
}

// Build will provide methods for interaction with API regarding building images
type Build struct {
	Client dockerlib.Client
}

/*
BuildImage builds a Docker image via the Docker API. Takes the source directory
and image name and then builds the appropriate image. Tarball is utilized
in order to make building easier.
*/
func (c *Build) BuildImage(dockerfile, image, context string) error {

	log.Infof("Building image '%s' from directory '%s'", image, path.Base(context))

	// Create a temporary file for tarball image packaging
	tmpFile, err := ioutil.TempFile("/tmp", "kedge-image-build-")
	// Delete tarball after creating image
	defer os.Remove(tmpFile.Name())

	if err != nil {
		return err
	}
	log.Debugf("Created temporary file %v for Docker image tarballing", tmpFile.Name())

	// Create a tarball of the source directory in order to build the resulting image
	err = CreateTarball(strings.Join([]string{context, ""}, "/"), tmpFile.Name())
	if err != nil {
		return errors.Wrap(err, "unable to create a tarball")
	}

	// Load the file into memory
	tarballSource, err := os.Open(tmpFile.Name())
	if err != nil {
		return errors.Wrap(err, "unable to load tarball into memory")
	}

	// Let's create all the options for the image building.
	outputBuffer := bytes.NewBuffer(nil)
	opts := dockerlib.BuildImageOptions{
		Name:         image,
		InputStream:  tarballSource,
		OutputStream: outputBuffer,
		Dockerfile:   dockerfile,
	}
	// Build it!
	if err := c.Client.BuildImage(opts); err != nil {
		return errors.Wrap(err, "unable to build image")
	}

	log.Infof("Image '%s' from directory '%s' built successfully", image, path.Base(context))
	log.Debugf("Image %s build output:\n%s", image, outputBuffer)
	return nil
}

/*
CreateTarball creates a tarball for source and dumps it to target path

Function modified and added from https://github.com/mholt/archiver/blob/master/tar.go
*/
func CreateTarball(source, target string) error {
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if baseDir == path {
				return nil
			}
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				if strings.HasSuffix(source, "/") {
					header.Name = strings.TrimPrefix(path, source)
				} else {
					header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
				}
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}
