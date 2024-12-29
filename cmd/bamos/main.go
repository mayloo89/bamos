package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/driver"
	"github.com/mayloo89/bamos/internal/handler"
	"github.com/mayloo89/bamos/internal/helpers"
	"github.com/mayloo89/bamos/internal/model"
	"github.com/mayloo89/bamos/internal/render"
	"github.com/mayloo89/bamos/utils"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	// Close the database connection when the main function returns.
	defer db.SQL.Close()

	fmt.Printf("starting application at port %s \n", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// What am I going to put in the config?
	gob.Register(model.Route{})

	// change this when in production
	app.InProduction = false

	// set up the loggers
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// Connect to database.
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bamos user=ssourigues")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("can not create template cache: " + err.Error())
	}

	app.TemplateCache = tc
	app.UseCache = false
	app.DataCache.Routes = utils.GetRoutes()
	if len(app.DataCache.Routes) <= 0 {
		fmt.Printf("no routes were loaded in the cache\n")
	}
	fmt.Printf("routes cache loaded with %d routes.\n", len(app.DataCache.Routes))

	repo := handler.NewRepo(&app, db)
	handler.NewHandler(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
