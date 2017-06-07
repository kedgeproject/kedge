package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCommand() *cobra.Command {
	v := viper.New()
	v.SetEnvPrefix("opencomposition")
	v.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	v.SetEnvKeyReplacer(replacer)

	var rootCmd = &cobra.Command{
		Use: "opencomposition",
		RunE: func(cmd *cobra.Command, args []string) error {
			//if err := cmd.Help(); err != nil {
			//	return err
			//}
			return errors.New("Use 'opencomposition convert'")
		},
	}

	rootCmd.AddCommand(NewConvertCommand(v))

	return rootCmd
}
