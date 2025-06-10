package handler

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"

	"github.com/mayloo89/bamos/internal/config"
)

var app config.AppConfig
var session *scs.SessionManager

// NoSurf adds CSRF protection to all POST request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session in every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	dir, err := filepath.Abs("./../../templates")
	if err != nil {
		return myCache, err
	}

	dir = strings.Replace(dir, "/cmd/bamos/", "/", 1)

	pagesGlob := dir + "/*.page.tmpl"
	layoutGlob := dir + "/*.layout.tmpl"

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

func TestCreateTestTemplateCache(t *testing.T) {
	// your test code here
}
