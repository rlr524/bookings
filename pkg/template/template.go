/*
Package template implements a simple means by which to create a template cache and render templates.
*/
package template

import (
	"bytes"
	"fmt"
	"github.com/rlr524/bookings/pkg/config"
	"github.com/rlr524/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData allows to use the TemplateData struct to add default data for all pages
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	// Add data here as needed
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// Create a new template cache
		tc, _ = CreateTemplateCache()
	}

	// Check that there is a template in the template cache map and assign it to t if so, else fail
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	// Create a buffer and assign it to buf
	buf := new(bytes.Buffer)

	// Pass in whatever data is added in the AddDefaultData function
	td = AddDefaultData(td)

	// Take the template, execute it, don't add any data, and store the value in the buffer
	_ = t.Execute(buf, td)
	// Write to the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}
}

// CreateTemplateCache creates a template cache as a map of template pages included in the templates directory
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.gohtml")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Create the template set (ts) from the page templates (*.page.gohtml)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
