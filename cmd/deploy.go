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
	pkgcmd "github.com/kedgeproject/kedge/pkg/cmd"
	"github.com/spf13/cobra"
	"os"
)

// Variables
var (
	DeployFiles []string
)

// Represents the "deploy" command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an application to Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if err := pkgcmd.Deploy(DeployFiles); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	deployCmd.Flags().StringArrayVarP(&DeployFiles, "files", "f", []string{}, "Specify files")
	RootCmd.AddCommand(deployCmd)
}
