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

// Represents the "delete" command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the resource from the Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ifFilesPassed(InputFiles); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		kubectlCommand := []string{"delete"}

		// Only setting the namespace flag to kubectl when --namespace is passed
		// explicitly by the user
		if cmd.Flags().Lookup("namespace").Changed {
			kubectlCommand = append(kubectlCommand, "--namespace", Namespace)
		}

		if err := pkgcmd.CreateKubernetesArtifacts(InputFiles, false, kubectlCommand...); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	deleteCmd.Flags().StringArrayVarP(&InputFiles, "files", "f", []string{}, "Specify files")
	deleteCmd.MarkFlagRequired("files")
	deleteCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "Kubernetes namespace to delete your application from")
	RootCmd.AddCommand(deleteCmd)
}
