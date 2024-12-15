package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/mayloo89/bamos/utils"
)

// AppConfig holds the application configurations
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	DataCache     struct {
		Routes []utils.Route
	}
}
