package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	pkgcmd "github.com/surajssd/kapp/pkg/cmd"
)

// Variables
var (
	GenerateFiles []string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Kubernetes resources from application definition(s)",
	Run: func(cmd *cobra.Command, args []string) {
		if err := pkgcmd.Generate(GenerateFiles); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	generateCmd.Flags().StringArrayVarP(&GenerateFiles, "files", "f", []string{}, "Specify files")
	RootCmd.AddCommand(generateCmd)
}
