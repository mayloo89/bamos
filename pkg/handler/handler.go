package handler

import (
	"net/http"

	"github.com/mayloo89/bamos/pkg/config"
	"github.com/mayloo89/bamos/pkg/model"
	"github.com/mayloo89/bamos/pkg/render"
)

var Repo *Repository

type (
	Repository struct {
		App *config.AppConfig
	}
)

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again."

	render.RenderTemplate(w, "home.page.tmpl", &model.TemplateData{
		StringMap: stringMap,
	})
}
