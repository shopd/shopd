package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/config"
	"github.com/shopd/shopd/go/router"
	"github.com/spf13/cobra"
)

type RunHandler struct {
	http.Server
	cleanup  func() error
	done     chan struct{}
	shutdown chan os.Signal
}

func NewRunHandler() *RunHandler {
	rh := &RunHandler{
		cleanup:  func() error { return nil },
		done:     make(chan struct{}),
		shutdown: make(chan os.Signal, 1),
	}
	return rh
}

type NewServerParams struct {
	Stubs bool
}

func NewServer(conf *config.Config, params NewServerParams) (
	rh *RunHandler, err error) {

	// TODO Services
	// s := services.NewServices(conf)

	if params.Stubs {
		log.Info().Msg("stubs")
		// TODO Refactor how stubs work
		// - ApiServer considers using stubs if this mode is set
		// - use stubs if route is not found
		// - use stubs if request header is set
		// - create empty stubs if not found and header is set
		// - otherwise use route handler
	}

	// Setup HTTP server
	r := router.NewRouter(conf)
	rh = NewRunHandler()
	rh.Server = http.Server{}
	rh.Handler = r.Handler()
	rh.Addr = conf.PortApi()
	// rh.cleanup = s.Cleanup
	rh.cleanup = func() error {
		log.Info().Msg("TODO Cleanup services")
		return nil
	}

	return rh, nil
}

func (rh *RunHandler) Run() (err error) {
	go func() {
		// shutdown chan initiates shutdown on interrupt signal
		// "Gracefully shut down the server
		// without interrupting active connections...
		// does not attempt to close nor wait for WebSockets"
		// https://golang.org/pkg/net/http/#Server.Shutdown
		signal.Notify(rh.shutdown, os.Interrupt, syscall.SIGTERM)
		<-rh.shutdown
		log.Info().Msg("ctrl+c interrupt, shutting down...")
		defer close(rh.done)

		// Interrupt signal received
		err := rh.Shutdown(context.Background())
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			return
		}

		// Cleanup
		rh.cleanup()
	}()

	// Listen to requests
	log.Info().Msgf("listening on %s", rh.Addr)
	err = rh.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			// Unexpected shutdown
			return errors.WithStack(err)
		}
	}
	// Graceful shutdown
	rh.wait()
	log.Info().Msg("bye!")

	return nil
}

// Signal on the shutdown channel,
// and wait for graceful shutdown to complete.
// Can be used instead of Ctrl+C when running the server from code,
// e.g. from mage task. Was used previously with the WatchServer target
func (rh *RunHandler) Signal() {
	// "SIGINT is usually user-initiated,
	// while SIGTERM can be system or process-initiated"
	rh.shutdown <- os.Signal(syscall.SIGTERM)
	rh.wait()
}

// wait for shutdown and cleanup to complete
func (rh *RunHandler) wait() {
	<-rh.done
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the shopd process and blocks indefinitely",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		conf := cmd.Context().Value(config.Config{}).(*config.Config)

		stubs, err := cmd.Flags().GetBool("stubs")
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}

		rh, err := NewServer(conf, NewServerParams{Stubs: stubs})
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			rh.cleanup()
			os.Exit(1)
		}

		// Run until shutdown signal is received
		err = rh.Run()
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
