package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/config"
	"github.com/shopd/shopd/go/fileutil"
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

		// API helper
		log.Info().Str("env", env).
			Strs("apiPaths", apiPaths).Msg("Generating api helper")
		t, err := template.New("apiTemplate").Parse(apiTemplate)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		buf := bytes.Buffer{}
		err = t.Execute(&buf, map[string]any{
			"Paths": apiPaths,
		})
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		fileutil.WriteBytes(
			filepath.Join(conf.Dir(), "go", "templ", "api_templ.go"),
			buf.Bytes())

		// Static site helper
		if env == share.EnvDev {
			log.Info().Str("env", env).
				Strs("contentPaths", contentPaths).
				Msg("Generating static site helper")
			t, err := template.New("contentTemplate").Parse(contentTemplate)
			if err != nil {
				log.Error().Stack().Err(err).Msg("")
				os.Exit(1)
			}
			buf := bytes.Buffer{}
			err = t.Execute(&buf, map[string]any{
				"Paths": contentPaths,
			})
			if err != nil {
				log.Error().Stack().Err(err).Msg("")
				os.Exit(1)
			}
			fileutil.WriteBytes(
				filepath.Join(conf.Dir(), "go", "templ", "static_templ.go"),
				buf.Bytes())

		} else {
			// Static site
			log.Info().Str("env", env).
				Strs("contentPaths", contentPaths).
				Msg("Generating static site")
			// TODO Generate www/public.
			// Copy contents of www/static to www/public.
			// Copy additional domain specific static content and overrides
		}

		os.Exit(0)
	},
}

func init() {
	staticCmd.AddCommand(staticGenCmd)
	staticGenCmd.Flags().String(FlagEnv, "", "")
	staticGenCmd.MarkFlagRequired(FlagEnv)
}

const apiTemplate = `
package templ
{{range .Paths}}
// {{.}}
{{end}}
`

const contentTemplate = `
package templ
{{range .Paths}}
// {{.}}
{{end}}
`
