package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/surajssd/opencomposition/pkg"
)

func NewConvertCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use: "convert",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunConvert(v, cmd)
		},
	}
	cmd.PersistentFlags().StringSliceP("files", "f", []string{}, "Specify opencompose files")
	v.BindPFlag("files", cmd.PersistentFlags().Lookup("files"))

	return cmd
}

func RunConvert(v *viper.Viper, cmd *cobra.Command) error {
	return errors.Wrap(pkg.Convert(v, cmd), "failed conversion")
}
