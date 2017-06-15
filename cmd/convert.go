package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	pkgcmd "github.com/surajssd/opencomposition/pkg/cmd"
)

// Variables
var (
	ConvertFiles []string
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert an application to Kubernetes resources",
	Run: func(cmd *cobra.Command, args []string) {
		if err := pkgcmd.Convert(ConvertFiles); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	convertCmd.Flags().StringArrayVarP(&ConvertFiles, "files", "f", []string{}, "Specify files")
	RootCmd.AddCommand(convertCmd)
}
