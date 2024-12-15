package handler

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"

	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/render"
	"github.com/mayloo89/bamos/utils"
)

var app config.AppConfig
var session *scs.SessionManager

func getRoutes() http.Handler {
	// change this when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("can not create template cache: " + err.Error())
	}

	app.TemplateCache = tc
	app.UseCache = true
	// set up the loggers
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.DataCache.Routes = utils.GetRoutes()
	if len(app.DataCache.Routes) <= 0 {
		fmt.Printf("no routes were loaded in the cache\n")
	}
	fmt.Printf("routes cache loaded with %d routes.\n", len(app.DataCache.Routes))

	repo := NewRepo(&app)
	NewHandler(repo)

	render.NewTemplates(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", Repo.Home)
	mux.Get("/colectivos/vehiclePositionsSimple", Repo.VehiclePositionsSimple)
	mux.Get("/colectivos/feed-gtfs-frequency", Repo.FeedGtfsFrequency)

	mux.Get("/colectivos/search", Repo.SearchLine)
	mux.Post("/colectivos/search", Repo.PostSearchLine)

	return mux
}

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
