package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	pkgcmd "github.com/surajssd/kapp/pkg/cmd"
)

// Variables
var (
	DeployFiles []string
)

// convertCmd represents the convert command
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
