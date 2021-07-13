package models

import "github.com/liomazza/bookings/internal/forms"

//TemplateData holds data sent from handlers to template
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]string
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string //Cross site Request Forgery Token
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
