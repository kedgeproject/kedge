package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Global variables
var (
	GlobalVerbose bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "opencomposition",
	Short: "Compose Kubernetes applications using Kubernetes constructs",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// Add extra logging when verbosity is passed
		if GlobalVerbose {
			log.SetLevel(log.DebugLevel)
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
}
