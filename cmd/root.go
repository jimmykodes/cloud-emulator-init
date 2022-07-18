package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jimmykodes/cloud-emulator-init/internal/runner"
)

var rootCmd = &cobra.Command{
	Use:          "cloud-emulator-init",
	Short:        "initialize cloud emulation resources",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: runner.RunE,
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd.Flags().AddFlagSet(runner.Flags())
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
