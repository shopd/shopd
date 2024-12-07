package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/config"
	"github.com/shopd/shopd/go/share"
	"github.com/spf13/cobra"
)

// findFilePaths matches a pattern in dir and returns a list of matching paths
func findFilePaths(dir, pattern string) (matches []string, err error) {
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
	Short: "Generate api helper and static site files",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		conf := cmd.Context().Value(config.Config{}).(*config.Config)

		env, err := cmd.Flags().GetString(FlagEnv)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}

		// Scan contents of www
		apiPaths := make([]string, 0)
		contentPaths := make([]string, 0)
		dir := filepath.Join(conf.Dir(), "www")
		pattern := "*_templ.go"
		wwwMatches, err := findFilePaths(dir, pattern)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		for _, match := range wwwMatches {
			if strings.Contains(match, "www/api") {
				apiPaths = append(apiPaths, match)
			} else if strings.Contains(match, "www/content") {
				contentPaths = append(contentPaths, match)
			}
		}

		log.Info().Str("env", env).Msg("Generating api helper")
		//  TODO Generate go/templ/api_templ.go
		// for paths starting with "/api"
		log.Info().Strs("apiPaths", apiPaths).Msg("")

		if env == share.EnvDev {
			log.Info().Str("env", env).Msg("Generating static site helper")
			//  TODO Generate www/static_gen.go dev service
			// for paths starting with "/content".
			// Caddy forwards static site requests to this service
			log.Info().Strs("contentPaths", contentPaths).Msg("")

		} else {
			log.Info().Str("env", env).Msg("Generating static site")
			// TODO Generate static site files in www/public.
			// Copy contents of www/static to www/public
		}

		os.Exit(0)
	},
}

func init() {
	staticCmd.AddCommand(staticGenCmd)
	staticGenCmd.Flags().String(FlagEnv, "", "")
	staticGenCmd.MarkFlagRequired(FlagEnv)
}
