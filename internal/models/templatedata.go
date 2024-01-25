package models

import "github.com/Soyaib10/comfort-cocoon/internal/forms"

// TemplateData holds data send from handlers to template
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	From      *forms.Form
}
