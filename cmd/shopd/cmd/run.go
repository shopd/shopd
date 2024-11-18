package cmd

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the shopd process and blocks indefinitely",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// conf := cmd.Context().Value(config.Config{}).(*config.Config)

		stubs, err := cmd.Flags().GetBool("stubs")
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}

		if stubs {
			log.Info().Msg("stubs")
			// TODO Refactor how stubs work
			// - ApiServer considers using stubs if this mode is set
			// - use stubs if route is not found
			// - use stubs if request header is set
			// - create empty stubs if not found and header is set
			// - otherwise use route handler
		}
		r := gin.Default()
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		// TODO Graceful shutdown
		// https://g.co/gemini/share/cb8bcb4a6b76
		// TODO Templ integration
		// https://templ.guide/integrations/web-frameworks/
		// https://github.com/a-h/templ/blob/main/examples/integration-gin/main.go

		// Run until shutdown signal is received
		err = r.Run()
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().Bool("stubs", false, "Enable stubs mode")
}
