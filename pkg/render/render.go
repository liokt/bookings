package render

import (
	"bytes"
	"fmt"
	"github.com/liomazza/bookings/pkg/config"
	"github.com/liomazza/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{

}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig){
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html templates
func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {

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

	templateData = AddDefaultData(templateData)

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

		//fmt.Println("Page is currently", page)

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