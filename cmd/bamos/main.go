package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mayloo89/bamos/pkg/config"
	"github.com/mayloo89/bamos/pkg/handler"
	"github.com/mayloo89/bamos/pkg/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	//TODO: change this when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can not create template cache: " + err.Error())
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handler.NewRepo(&app)
	handler.NewHandler(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application at port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
