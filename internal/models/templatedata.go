package models

import "github.com/ol-ilyassov/booking-app/internal/forms"

// TemplateData holds data sent from handlers to templates.
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float64
	Data      map[string]interface{}
	// In terms of forms work and (post) methods processing in handlers,
	// a security token = cross-site request forgery.
	CSRFToken string
	// Response messages:
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}
