package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// VERSION  is version number that will be displayed when running ./kedge version
	VERSION = "0.0.0"

	// GITCOMMIT is hash of the commit that wil be displayed when running ./kedge version
	// this will be overwritten when running  build like this: go build -ldflags="-X github.com/kedgeproject/kedge/cmd.GITCOMMIT=$(GITCOMMIT)"
	// HEAD is default indicating that this was not set during build
	GITCOMMIT = "HEAD"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Kedge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION + " (" + GITCOMMIT + ")")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
