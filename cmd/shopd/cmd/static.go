package cmd

import (
	"github.com/spf13/cobra"
)

// staticCmd represents the static command
var staticCmd = &cobra.Command{
	Use:   "static",
	Short: "Static site management",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(staticCmd)
}
