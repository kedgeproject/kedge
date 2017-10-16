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

	pkgcmd "github.com/kedgeproject/kedge/pkg/cmd"
	"github.com/spf13/cobra"
)

// Represents the "create" command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the resource on the Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ifFilesPassed(InputFiles); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		command := []string{"create", "--save-config"}

		// Only setting the namespace flag to kubectl when --namespace is passed
		// explicitly by the user
		if cmd.Flags().Lookup("namespace").Changed {
			command = append(command, "--namespace", Namespace)
		}

		if err := pkgcmd.CreateArtifacts(InputFiles, false, command...); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	createCmd.Flags().StringArrayVarP(&InputFiles, "files", "f", []string{}, "Specify files")
	createCmd.MarkFlagRequired("files")
	createCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "Namespace or project to deploy your application to")

	RootCmd.AddCommand(createCmd)
}
