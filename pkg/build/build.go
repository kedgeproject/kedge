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
	"github.com/ghodss/yaml"
	os_build_v1 "github.com/openshift/origin/pkg/build/apis/build/v1"
	os_image_v1 "github.com/openshift/origin/pkg/image/apis/image/v1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/kubernetes/pkg/api/install"
	kapi "k8s.io/kubernetes/pkg/api/v1"

	"github.com/kedgeproject/kedge/pkg/cmd"
	"github.com/kedgeproject/kedge/pkg/spec"
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
	log.Debugf("Image %s build output:\n%s", image, spec.PrettyPrintObjects(outputBuffer))
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
			var symLink string
			// Checking for the symlinks
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				// Allow symlinks
				if symLink, err = os.Readlink(path); err != nil {
					return err
				}
			}

			fileHeader, err := tar.FileInfoHeader(info, symLink)
			if err != nil {
				return err
			}

			if baseDir != "" {
				if strings.HasSuffix(source, "/") {
					fileHeader.Name = strings.TrimPrefix(path, source)
				} else {
					fileHeader.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
				}
			}

			if err := tarball.WriteHeader(fileHeader); err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
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

// getImageTag get tag name from image name
// if no tag is specified return 'latest'
func GetImageTag(image string) string {
	// format:      registry_host:registry_port/repo_name/image_name:image_tag
	// example:
	// 1)     myregistryhost:5000/fedora/httpd:version1.0
	// 2)     myregistryhost:5000/fedora/httpd
	// 3)     myregistryhost/fedora/httpd:version1.0
	// 4)     myregistryhost/fedora/httpd
	// 5)     fedora/httpd
	// 6)     httpd
	imageAndTag := image

	imageTagSplit := strings.Split(image, "/")
	if len(imageTagSplit) >= 2 {
		imageAndTag = imageTagSplit[len(imageTagSplit)-1]
	}

	p := strings.Split(imageAndTag, ":")
	if len(p) == 2 {
		return p[1]
	}
	return "latest"
}

func GetImageName(image string) string {
	imageAndTag := image

	imageTagSplit := strings.Split(image, "/")
	if len(imageTagSplit) >= 2 {
		imageAndTag = imageTagSplit[len(imageTagSplit)-1]
	}
	p := strings.Split(imageAndTag, ":")
	if len(p) <= 2 {
		return p[0]
	}

	return image
}

func BuildS2I(image, context, builderImage string, namespace string) error {

	name := GetImageName(image)
	labels := map[string]string{
		spec.BuildLabelKey: name,
	}
	annotations := map[string]string{
		"openshift.io/generated-by": "KedgeBuildS2I",
	}

	is := os_image_v1.ImageStream{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ImageStream",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: os_image_v1.ImageStreamSpec{
			Tags: []os_image_v1.TagReference{
				{Name: GetImageTag(image)},
			},
		},
	}

	bc := os_build_v1.BuildConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BuildConfig",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: os_build_v1.BuildConfigSpec{
			CommonSpec: os_build_v1.CommonSpec{
				Strategy: os_build_v1.BuildStrategy{
					Type: "Binary",
					SourceStrategy: &os_build_v1.SourceBuildStrategy{
						From: kapi.ObjectReference{
							Kind: "DockerImage",
							Name: builderImage,
						},
					},
				},
				Output: os_build_v1.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "ImageStreamTag",
						Name: name + ":" + GetImageTag(image),
					},
				},
			},
		},
	}

	isyaml, err := yaml.Marshal(is)
	if err != nil {
		return err
	}

	bcyaml, err := yaml.Marshal(bc)
	if err != nil {
		return err
	}

	log.Debugf("ImageStream for output image: \n%s\n", string(isyaml))
	log.Debugf("BuildConfig: \n%s\n", string(bcyaml))

	args := []string{"apply", "-f", "-"}
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}
	err = cmd.RunClusterCommand(args, isyaml, true)
	if err != nil {
		return err
	}
	err = cmd.RunClusterCommand(args, bcyaml, true)
	if err != nil {
		cleanup(name)
		return err
	}

	log.Infof("Starting build for %q", image)
	cmd := []string{"oc", "start-build", image, "--from-dir=" + context, "-F", "--namespace", namespace}
	if err := RunCommand(cmd); err != nil {
		return err
	}

	return nil
}

func cleanup(name string) {
	log.Infof("Cleaning up build since error occurred while building")

	delBc := []string{"oc", "delete", "buildconfig", name}
	if err := RunCommand(delBc); err != nil {
		log.Debugf("error while deleting buildconfig: %v", err)
	}

	delIs := []string{"oc", "delete", "imagestream", name}
	if err := RunCommand(delIs); err != nil {
		log.Debugf("error while deleting imagestream: %v", err)
	}
}
