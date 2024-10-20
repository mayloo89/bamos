package render

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
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

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, tmplData *model.TemplateData) {
	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		// get template cache from app config
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Template " + tmpl + " does not exist in cache (len " + strconv.Itoa(len(tc)) + ")")
	}

	buf := new(bytes.Buffer)

	tmplData = AddDefaultData(tmplData, r)

	err = t.Execute(buf, tmplData)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	dir, err := filepath.Abs("./templates/")
	if err != nil {
		return myCache, err
	}

	pagesGlob := filepath.Join(dir, "*.page.tmpl")
	layoutGlob := filepath.Join(dir, "*.layout.tmpl")

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
