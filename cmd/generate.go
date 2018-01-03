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

// Represents the "generate" command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Kubernetes resources from an app definition",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ifFilesPassed(InputFiles); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		if err := pkgcmd.CreateArtifacts(InputFiles, true, SkipValidation, ""); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	generateCmd.Flags().StringArrayVarP(&InputFiles, "files", "f", []string{}, "Input files")
	generateCmd.MarkFlagRequired("files")
	generateCmd.Flags().BoolVar(&SkipValidation, "skip-validation", false, "Skip validation of Kedge file")
	RootCmd.AddCommand(generateCmd)
}
