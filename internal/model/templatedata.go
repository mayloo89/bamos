package model

import (
	"html/template"

	"github.com/mayloo89/bamos/internal/forms"
)

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap   map[string]string
	IntMap      map[string]int
	FloatMap    map[string]float32
	Data        map[string]interface{}
	TemplateMap map[string]template.HTML
	CSRFToken   string
	Flash       string
	Warning     string
	Error       string
	Form        *forms.Form
}
