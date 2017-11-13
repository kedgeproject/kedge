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
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	dockerlib "github.com/fsouza/go-dockerclient"
	"github.com/ghodss/yaml"
	os_build_v1 "github.com/openshift/origin/pkg/build/apis/build/v1"
	os_image_v1 "github.com/openshift/origin/pkg/image/apis/image/v1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kapi "k8s.io/kubernetes/pkg/api/v1"

	"github.com/kedgeproject/kedge/pkg/cmd"
	"github.com/kedgeproject/kedge/pkg/spec"
	_ "k8s.io/kubernetes/pkg/api/install"
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

// GetGitCurrentBranch gets current git branch name for the current git repo
func GetGitCurrentBranch(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil
}

// GetGitCurrentRemoteURL gets current git remote URI for the current git repo
func GetGitCurrentRemoteURL(dir string) (string, error) {
	cmd := exec.Command("git", "ls-remote", "--get-url")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	url := strings.TrimRight(string(out), "\n")
	if !strings.HasSuffix(url, ".git") {
		url += ".git"
	}

	if strings.HasPrefix(url, "git") {
		// URL: git@github.com:surajssd/kedge.git
		// will be divided into tokens like this:
		// git@ github.com surajssd/kedge.git
		// src: https://stackoverflow.com/a/2514986/3848679
		// More generic regex: '(\w+://)(.+@)*([\w\d\.]+)(:[\d]+){0,1}/*(.*)'
		urlRe := regexp.MustCompile(`(.+@)*([\w\d\.]+):(.*)`)
		urlComponents := urlRe.FindStringSubmatch(url)
		url = filepath.Join(urlComponents[len(urlComponents)-2],
			urlComponents[len(urlComponents)-1])
		url = "https://" + url
	}
	return url, nil
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

func localdir() (string, error) {
	abs, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	return filepath.Base(abs), nil
}

func BuildS2I(dockerfile, image, context string) error {
	repo, err := GetGitCurrentRemoteURL(context)
	if err != nil {
		return err
	}
	branch, err := GetGitCurrentBranch(context)
	if err != nil {
		return err
	}
	// name of this build
	//name, err := localdir()
	//if err != nil {
	//	return err
	//}
	name := GetImageName(image)
	// labels
	labels := spec.GetNameLabel(name)

	bc := os_build_v1.BuildConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BuildConfig",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: os_build_v1.BuildConfigSpec{
			Triggers: []os_build_v1.BuildTriggerPolicy{
				{Type: "ConfigChange"},
			},
			CommonSpec: os_build_v1.CommonSpec{
				Source: os_build_v1.BuildSource{
					Git: &os_build_v1.GitBuildSource{
						URI: repo,
						Ref: branch,
					},
					ContextDir: context,
					Type:       "Git",
				},
				Strategy: os_build_v1.BuildStrategy{
					Type: "Docker",
					DockerStrategy: &os_build_v1.DockerBuildStrategy{
						DockerfilePath: dockerfile,
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

	is := os_image_v1.ImageStream{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ImageStream",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: os_image_v1.ImageStreamSpec{
			LookupPolicy: os_image_v1.ImageLookupPolicy{
				Local: false,
			},
			Tags: []os_image_v1.TagReference{
				{
					From: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: image,
					},
				},
			},
		},
	}

	bcyaml, err := yaml.Marshal(bc)
	if err != nil {
		return err
	}
	isyaml, err := yaml.Marshal(is)
	if err != nil {
		return err
	}
	args := []string{"create", "-f", "-"}
	err = cmd.RunClusterCommand(args, isyaml, true)
	if err != nil {
		return err
	}
	err = cmd.RunClusterCommand(args, bcyaml, true)
	if err != nil {
		return err
	}
	//fmt.Println("---")
	//fmt.Println(string(bcyaml))
	//fmt.Println("---")
	//fmt.Println(string(isyaml))

	return nil
}
