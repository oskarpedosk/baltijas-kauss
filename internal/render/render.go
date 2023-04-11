package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/oskarpedosk/baltijas-kauss/internal/config"
	"github.com/oskarpedosk/baltijas-kauss/internal/models"
)

var pathToTemplates = "./templates"
var functions = template.FuncMap{
	"add": add,
	"seq": seq,
}

var app *config.AppConfig

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(tmplData *models.TemplateData, r *http.Request) *models.TemplateData {
	tmplData.Flash = app.Session.PopString(r.Context(), "flash")
	tmplData.Error = app.Session.PopString(r.Context(), "error")
	tmplData.Info = app.Session.PopString(r.Context(), "info")
	tmplData.Warning = app.Session.PopString(r.Context(), "warning")
	tmplData.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		tmplData.User = models.User{
			UserID:      app.Session.GetInt(r.Context(), "user_id"),
			FirstName:   app.Session.GetString(r.Context(), "first_name"),
			LastName:    app.Session.GetString(r.Context(), "last_name"),
			Email:       app.Session.GetString(r.Context(), "email"),
			ImgID:       app.Session.GetString(r.Context(), "img"),
			AccessLevel: app.Session.GetInt(r.Context(), "access_level"),
		}
	}
	return tmplData
}

// RenderTemplate renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, templateName string, tmplData *models.TemplateData) error {
	var templateCache map[string]*template.Template

	if app.UseCache {
		// Get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		t, err := CreateTemplateCache()
		if err != nil {
			log.Println(err)
		}
		templateCache = t
	}

	tmpl, ok := templateCache[templateName]
	if !ok {
		return errors.New("can't get template from cache")
	}

	tmplData = AddDefaultData(tmplData, r)

	buf := new(bytes.Buffer)

	err := tmpl.Execute(buf, tmplData)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error executing template data")
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all the files named *.page.tmpl from ../../templates/
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// Range through all the files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}

	return myCache, nil
}

func add(x, y int) int {
	return x + y
}

func seq(start, end int) []int {
	var s []int
	for i := start; i <= end; i++ {
		s = append(s, i)
	}
	return s
}
