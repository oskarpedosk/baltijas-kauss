package models

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap  map[string]string
	IntMap     map[string]int
	PlayerData map[string]interface{}
	Data       map[string]interface{}
	CSRFToken  string
	Flash      string
	Warning    string
	Error      string
}
