package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var staticDevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Dev service for static site",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// conf := cmd.Context().Value(config.Config{}).(*config.Config)

		// TODO See comments for Caddy in Makefile

		os.Exit(0)
	},
}

func init() {
	staticCmd.AddCommand(staticDevCmd)
}
