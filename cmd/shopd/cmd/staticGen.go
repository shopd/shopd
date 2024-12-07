package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/config"
	"github.com/spf13/cobra"
)

// staticGenMatch find matches for pattern in dir
// and returns a list of matching paths
func staticGenMatch(dir, pattern string) (matches []string, err error) {
	err = filepath.WalkDir(dir,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				return err
			} else if matched {
				matches = append(matches, path)
			}

			return nil
		})

	if err != nil {
		return matches, err
	}

	return matches, nil
}

var staticGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate static site or helpers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		conf := cmd.Context().Value(config.Config{}).(*config.Config)

		env, err := cmd.Flags().GetString(FlagEnv)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}

		// TODO Generate static site or helpers.
		// WatchShopd task watches for changes to the Go source,
		// and rebuilds the shopd backend.
		// On dev the helpers are used to serve static site requests,
		// except for the files in www/static
		log.Info().Str("env", env).Msg("staticGenCmd")

		// Scan contents of www
		dir := filepath.Join(conf.Dir(), "www")
		pattern := "*_templ.go"
		wwwMatches, err := staticGenMatch(dir, pattern)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		for _, match := range wwwMatches {
			fmt.Println(match)
		}

		// Scan go/templ
		dir = filepath.Join(conf.Dir(), "go", "templ")
		pattern = "*_gen.go"
		goMatches, err := staticGenMatch(dir, pattern)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		for _, match := range goMatches {
			fmt.Println(match)
		}

		// Compare matches

		os.Exit(0)
	},
}

func init() {
	staticCmd.AddCommand(staticGenCmd)
	staticGenCmd.Flags().String(FlagEnv, "", "")
	staticGenCmd.MarkFlagRequired(FlagEnv)
}
