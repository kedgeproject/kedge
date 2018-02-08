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

package cmd

import (
	"fmt"
	"os"

	"github.com/kedgeproject/kedge/pkg/build"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Dockerfile string
var DockerImage, BuilderImage string
var DockerContext string
var PushImage, s2iBuild bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build application image",
	Run: func(cmd *cobra.Command, args []string) {
		if DockerImage == "" {
			fmt.Println("Please specify the container image name using flag '--image' or '-i'")
			os.Exit(-1)
		}

		if s2iBuild {
			if PushImage {
				log.Warnf("Using source to image strategy for build, image will be by default pushed to internal container registry, so ignoring this flag")
			}
			if BuilderImage == "" {
				fmt.Println("Please specify the builder image name using flag '--builder-image' or '-b'")
				os.Exit(-1)
			}
			if err := build.BuildS2I(DockerImage, DockerContext, BuilderImage, Namespace); err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		} else {
			if err := build.BuildPushDockerImage(Dockerfile, DockerImage, DockerContext, PushImage); err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}
	},
}

func init() {

	buildCmd.Flags().StringVarP(&Dockerfile, "file", "f", "Dockerfile", "Specify Dockerfile for doing builds, Dockerfile path is relative to context")
	buildCmd.Flags().StringVarP(&DockerImage, "image", "i", "", "Image name and tag of resulting image")
	buildCmd.Flags().StringVarP(&DockerContext, "context", "c", ".", "Path to a directory containing a Dockerfile, it is build context that is sent to the Docker daemon")
	buildCmd.Flags().BoolVarP(&PushImage, "push", "p", false, "Add this flag if you want to push the image. Note: Ignored when s2i build strategy used")
	buildCmd.Flags().BoolVarP(&s2iBuild, "s2i", "", false, "If this is enabled then Source to Image build strategy is used")
	buildCmd.Flags().StringVarP(&BuilderImage, "builder-image", "b", "", "Name of a Docker image to use as a builder. Note: This is only useful when using s2i build strategy")
	buildCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "Namespace or project to deploy your application to")

	RootCmd.AddCommand(buildCmd)
}
