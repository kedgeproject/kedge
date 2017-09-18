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
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Global variables
var (
	GlobalVerbose  bool
	GlobalProvider string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kedge",
	Short: "Kedge: Simple, Concise & Declarative Kubernetes Applications",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// Add extra logging when verbosity is passed
		if GlobalVerbose {
			log.SetLevel(log.DebugLevel)
		}

		// Error out of the user has not chosen Kubernetes or OpenShift
		provider := strings.ToLower(GlobalProvider)
		if provider != "kubernetes" && provider != "openshift" {
			log.Fatalf("%s is an unsupported provider. Supported providers are: 'kubernetes', 'openshift'.", GlobalProvider)
		}

	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Initialize all flags
func init() {
	RootCmd.PersistentFlags().BoolVarP(&GlobalVerbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVar(&GlobalProvider, "provider", "kubernetes", "Specify a provider. Kubernetes or OpenShift.")
}
