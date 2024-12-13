package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is hard-coded here,
// and updated before creating new github releases
var version string = "v0.1.0"

func Version() string {
	return version
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
