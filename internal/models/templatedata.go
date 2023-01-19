package models

import "github.com/oskarpedosk/baltijas-kauss/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FuncMap         map[string]any
	Data            map[string]interface{}
	IntSlice        []int
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
	AccessLevel     int
}
