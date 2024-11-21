package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var staticGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate static site",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// conf := cmd.Context().Value(config.Config{}).(*config.Config)

		// TODO Generate static site to map www/content,
		// e.g. domains/example.com/www/public

		os.Exit(0)
	},
}

func init() {
	staticCmd.AddCommand(staticGenCmd)
}
