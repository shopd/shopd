package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
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

type templComponent struct {
	RelPath string
	Route   string
	Method  string
}

const apiPrefix = "/www/api"
const contentPrefix = "/www/content"

func stripBase(relPath string) string {
	return strings.Replace(relPath, filepath.Base(relPath), "", 1)
}

func stripTemplExt(relPath string) string {
	return strings.Replace(filepath.Base(relPath), "_templ.go", "", 1)
}

func apiMethod(segment string) (method string, err error) {
	switch segment {
	case "delete":
		return "Delete", nil
	case "get":
		return "Get", nil
	case "post":
		return "Post", nil
	case "put":
		return "Put", nil
	}
	return method, ErrNotImplemented(fmt.Sprintf("method %s", segment))
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
		apiComponents := make([]templComponent, 0)
		contentComponents := make([]templComponent, 0)
		dir := filepath.Join(conf.Dir(), "www")
		pattern := "*_templ.go"
		wwwMatches, err := findFilePaths(dir, pattern)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		for _, match := range wwwMatches {
			// See comments in /www/api and /www/content/README.md
			// on file naming conventions encoded in the logic below
			relPath := strings.Replace(match, conf.Dir(), "", 1)
			segment := stripTemplExt(relPath)

			if strings.Contains(relPath, apiPrefix) {
				method, err := apiMethod(segment)
				if err != nil {
					log.Error().Stack().Err(err).Msg("")
					os.Exit(1)
				}
				apiComponents = append(apiComponents, templComponent{
					RelPath: relPath,
					Route: path.Join("/", "api",
						strings.Replace(stripBase(relPath), apiPrefix, "", 1)),
					Method: method,
				})

			} else if strings.Contains(relPath, contentPrefix) {
				index := "index"
				route := strings.Replace(stripBase(relPath), contentPrefix, "", 1)
				if !strings.Contains(relPath, index) {
					route = path.Join(route, segment)
				}
				contentComponents = append(contentComponents, templComponent{
					RelPath: relPath,
					Route:   route,
					Method:  "Index",
				})
			}
		}

		// API helper
		log.Info().Str("env", env).Msg("Register api components")
		t, err := template.New("apiTemplate").Parse(apiTemplate)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		buf := bytes.Buffer{}
		err = t.Execute(&buf, map[string]any{
			"Components": apiComponents,
		})
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		fileutil.WriteBytes(
			filepath.Join(conf.Dir(), "go", "router", "api_templ.go"),
			buf.Bytes())

		// Static site helper
		if env == share.EnvDev {
			log.Info().Str("env", env).Msg("Register static site components")
			t, err := template.New("contentTemplate").Parse(contentTemplate)
			if err != nil {
				log.Error().Stack().Err(err).Msg("")
				os.Exit(1)
			}
			buf := bytes.Buffer{}
			err = t.Execute(&buf, map[string]any{
				"Components": contentComponents,
			})
			if err != nil {
				log.Error().Stack().Err(err).Msg("")
				os.Exit(1)
			}
			fileutil.WriteBytes(
				filepath.Join(conf.Dir(), "go", "router", "static_templ.go"),
				buf.Bytes())

		} else {
			// Static site
			log.Info().Str("env", env).Msg("Generating static site")
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
package router

{{range .Components}}
// TODO Import 
// {{.RelPath}}
// {{.Route}}
// {{.Method}}
{{end}}

func init() {
{{range .Components}}
// TODO Register component {{.}}
{{end}}
}
`

const contentTemplate = `
package router

{{range .Components}}
// TODO Import 
// {{.RelPath}}
// {{.Route}}
// {{.Method}}
{{end}}

func init() {
{{range .Components}}
// TODO Register component {{.}}
{{end}}
}
`
