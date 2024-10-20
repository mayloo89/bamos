package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
	"github.com/mayloo89/bamos/utils"
)

// AppConfig holds the application configurations
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
	DataCache     struct {
		Routes []utils.Route
	}
}
