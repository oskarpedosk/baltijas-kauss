package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

var pathToTemplates = "../../templates"
var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplate sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(tmplData *models.TemplateData, r *http.Request) *models.TemplateData {
	tmplData.Flash = app.Session.PopString(r.Context(), "flash")
	tmplData.Error = app.Session.PopString(r.Context(), "error")
	tmplData.Warning = app.Session.PopString(r.Context(), "warning")
	tmplData.CSRFToken = nosurf.Token(r)
	return tmplData
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, tmplData *models.TemplateData) error {

	var templateCache map[string]*template.Template
	if app.UseCache {
		// Get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}
 
	tmpl, ok := templateCache[templateName]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	tmplData = AddDefaultData(tmplData, r)

	buf := new(bytes.Buffer)

	err := tmpl.Execute(buf, tmplData)
	if err != nil {
		fmt.Println("Error executing template data")
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all the files named *.page.tmpl from ../../templates/
	pages, err := filepath.Glob("../../templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// Range through all the files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("../../templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("../../templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}

	return myCache, nil
}
