package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var staticGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate static site or helpers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// conf := cmd.Context().Value(config.Config{}).(*config.Config)

		env, err := cmd.Flags().GetString(FlagEnv)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}

		// TODO Generate static site or helpers.
		// WatchShopd task watches for changes on the helpers,
		// and rebuilds the shopd backend.
		// On dev the helpers are used to serve static site requests,
		// except for the files in www/static
		log.Info().Str("env", env).Msg("staticGenCmd")

		os.Exit(0)
	},
}

func init() {
	staticCmd.AddCommand(staticGenCmd)
	staticGenCmd.Flags().String(FlagEnv, "", "")
	staticGenCmd.MarkFlagRequired(FlagEnv)
}
