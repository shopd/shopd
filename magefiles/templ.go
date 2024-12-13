package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/fileutil"
	"github.com/shopd/shopd/go/share"
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

type templPackages struct {
	// m is a map of unique package paths
	m map[string]bool
}

func (tp *templPackages) Add(p string) {
	tp.m[p] = true
}

func (tp *templPackages) Packages() (packages []templPackage) {
	for v := range tp.m {
		packages = append(packages, templPackage{
			PackagePath: v,
		})
	}
	return packages
}

func NewTemplPackages() (tp *templPackages) {
	tp = &templPackages{}
	tp.m = make(map[string]bool)
	return tp
}

type templPackage struct {
	PackagePath string
}

type templComponent struct {
	FilePath    string
	PackageName string
	Constructor string
	Route       string
	Method      string
}

const apiPrefix = "/www/api"
const contentPrefix = "/www/content"

// stripBase removes the last element and any trailing slashes from the path
func stripBase(filePath string) string {
	return strings.TrimRight(
		strings.Replace(filePath, filepath.Base(filePath), "", 1), "/")
}

func stripTemplExt(filePath string) (packagePath, packageName, fileName string) {
	packagePath = filepath.Dir(filePath)
	packageName = filepath.Base(packagePath)
	fileName = strings.Replace(filepath.Base(filePath), "_templ.go", "", 1)
	return packagePath, packageName, fileName
}

// apiMethod maps filename to templ component constructors
func apiMethod(fileName string) (method string, err error) {
	switch strings.ToLower(fileName) {
	case "delete":
		return "Delete", nil
	case "get":
		return "Get", nil
	case "post":
		return "Post", nil
	case "put":
		return "Put", nil
	}
	return method, ErrNotImplemented(fmt.Sprintf("method %s", fileName))
}

// TemplApiGen generates router init file for templ api templates
func TemplApiGen(env string) (err error) {
	apiPackages := NewTemplPackages()

	// Scan contents of www
	apiComponents := make([]templComponent, 0)
	dir := filepath.Join(conf.Dir(), "www", "api")
	pattern := "*_templ.go"
	wwwMatches, err := findFilePaths(dir, pattern)
	if err != nil {
		return err
	}
	for _, match := range wwwMatches {
		filePath := strings.Replace(match, conf.Dir(), "", 1)
		packagePath, packageName, fileName := stripTemplExt(filePath)

		if strings.Contains(filePath, apiPrefix) {
			// api
			apiPackages.Add(packagePath)
			templConstructor, err := apiMethod(fileName)
			if err != nil {
				log.Error().Stack().Err(err).Msg("")
				os.Exit(1)
			}
			apiComponents = append(apiComponents, templComponent{
				FilePath:    filePath,
				PackageName: packageName,
				Constructor: templConstructor,
				Route: path.Join("/", "api",
					strings.Replace(stripBase(filePath), apiPrefix, "", 1)),
				Method: strings.ToUpper(templConstructor),
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
		"Packages":   apiPackages.Packages(),
		"Components": apiComponents,
	})
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		os.Exit(1)
	}
	fileutil.WriteBytes(
		filepath.Join(conf.Dir(), "go", "router", "init_api_templ.go"),
		buf.Bytes())

	return nil
}

// TemplSiteGen generates router init file for templ content templates
// TODO Generate static site files for prod?
func TemplSiteGen(env string) (err error) {
	contentPackages := NewTemplPackages()

	// Scan contents of www
	contentComponents := make([]templComponent, 0)
	dir := filepath.Join(conf.Dir(), "www", "content")
	pattern := "*_templ.go"
	wwwMatches, err := findFilePaths(dir, pattern)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		os.Exit(1)
	}
	for _, match := range wwwMatches {
		// See comments in /www/api and /www/content/README.md
		// on file naming conventions encoded in the logic below
		filePath := strings.Replace(match, conf.Dir(), "", 1)
		packagePath, packageName, fileName := stripTemplExt(filePath)

		if strings.Contains(filePath, contentPrefix) {
			// content
			contentPackages.Add(packagePath)
			index := "index"
			route := strings.Replace(stripBase(filePath), contentPrefix, "", 1)
			if !strings.Contains(filePath, index) {
				route = path.Join(route, fileName)
			}
			contentComponents = append(contentComponents, templComponent{
				FilePath:    filePath,
				PackageName: packageName,
				Route:       route,
				Constructor: "Index",
				Method:      share.GET,
			})
		}
	}

	// Static site helper
	if env == EnvDev {
		log.Info().Str("env", env).Msg("Register static site components")
		t, err := template.New("contentTemplate").Parse(contentTemplate)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		buf := bytes.Buffer{}
		err = t.Execute(&buf, map[string]any{
			"Packages":   contentPackages.Packages(),
			"Components": contentComponents,
		})
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		fileutil.WriteBytes(
			filepath.Join(conf.Dir(), "go", "router", "init_static_templ.go"),
			buf.Bytes())

	} else {
		// Static site
		return ErrNotImplemented(fmt.Sprintf("for env %s", env))
		// TODO Generate www/public.
		// Copy contents of www/static to www/public.
		// Copy additional domain specific static content and overrides
	}

	return nil
}

// TODO Split this into "api gen" command
const apiTemplate = `
package router

import ({{range .Packages}}
	"github.com/shopd/shopd{{.PackagePath}}"{{end}}
)

func init() {
	RegisterAPI = func(tr *TemplRegistry) {
		{{range .Components}}
		// {{.FilePath}}
		tr.RegisterAPI("{{.Route}}", "{{.Method}}", {{.PackageName}}.{{.Constructor}}()){{end}}
	}
}
`

const contentTemplate = `
package router

import ({{range .Packages}}
	"github.com/shopd/shopd{{.PackagePath}}"{{end}}
)

func init() {
{{range .Components}}
	RegisterStatic = func(tr *TemplRegistry) {
		// {{.FilePath}}
		tr.RegisterAPI("{{.Route}}", "{{.Method}}", {{.PackageName}}.{{.Constructor}}())
	}
{{end}}
}
`
