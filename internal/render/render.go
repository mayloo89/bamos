package render

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/justinas/nosurf"
	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/model"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(tmplData *model.TemplateData, r *http.Request) *model.TemplateData {
	tmplData.CSRFToken = nosurf.Token(r)
	return tmplData
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, tmplData *model.TemplateData) error {
	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		// get template cache from app config
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			// FIXME: return error here
			log.Fatal(err)
		}
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("Template " + tmpl + " does not exist in cache (len " + strconv.Itoa(len(tc)) + ")")
	}

	buf := new(bytes.Buffer)

	tmplData = AddDefaultData(tmplData, r)

	err = t.Execute(buf, tmplData)
	if err != nil {
		return err
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get the current working directory
	rootPath, err := getRootPath()
	if err != nil {
		return myCache, err
	}

	pagesGlob := rootPath + "/templates/*.page.tmpl"    // filepath.Join(dir, "*.page.tmpl")
	layoutGlob := rootPath + "/templates/*.layout.tmpl" // filepath.Join(dir, "*.layout.tmpl")

	// get all the *page.tmpl files from templates dir
	pages, err := filepath.Glob(pagesGlob)
	if err != nil {
		return myCache, err
	} else if len(pages) <= 0 {
		return myCache, errors.New("no pages found")
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		layouts, err := filepath.Glob(layoutGlob)
		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob(layoutGlob)
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

// getRootPath returns the root path of the application
func getRootPath() (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Find the root directory (assuming the root directory contains a specific file or directory)
	rootPath := cwd
	for {
		if _, err := os.Stat(filepath.Join(rootPath, "go.mod")); err == nil {
			break
		}

		parent := filepath.Dir(rootPath)
		if parent == rootPath {
			return "", errors.New("root path not found")
		}
		rootPath = parent
	}

	return rootPath, nil
}
