package render

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	config2 "github.com/liomazza/bookings/internal/config"
	models2 "github.com/liomazza/bookings/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{

}

var app *config2.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config2.AppConfig){
	app = a
}

func AddDefaultData(td *models2.TemplateData, r *http.Request) *models2.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, templateData *models2.TemplateData) {

	var myTemplateCache map[string]*template.Template

	if app.UseCache {
		//get the template cache from the app config
		myTemplateCache = app.TemplateCache
	} else {
		myTemplateCache, _ = CreateTemplateCache()
	}

	t, ok := myTemplateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buffer := new(bytes.Buffer)

	templateData = AddDefaultData(templateData, r)

	//We store the value of the template in a buffer so then we can read it
	_ = t.Execute(buffer, templateData)

	_, err := buffer.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

//CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() ( map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		fmt.Println("Page is currently", page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}