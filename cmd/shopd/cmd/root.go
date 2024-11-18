package cmd

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shopd",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// TODO Override context config for specified domain using config.LoadFile?
		domain, err := cmd.Flags().GetString(FlagDomain)
		if err != nil {
			log.Error().Stack().Err(errors.WithStack(err)).Msg("")
			os.Exit(1)
		}
		if domain != "" {
			err = errors.WithStack(ErrNotImplemented("domain flag"))
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
	},
}

func Execute(configBase64 string) {
	if configBase64 != "" {
		// Set compile time config
		err := config.SetEnvBase64(configBase64)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
	}

	// Set config on command context
	conf := config.New()
	ctx := context.Background()
	ctx = context.WithValue(ctx, config.Config{}, conf)

	// Execute root command
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().String(FlagDomain, "", "Specify domain")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
