package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
)

type templateData struct {
	Csrf       template.HTML
	Data       map[string]interface{}
	User       string
	FlashError string
	FlashInfo  string
}

// Render a template stack to html
// Pass the page template as the first template file
func (app *application) renderTemplates(w http.ResponseWriter, r *http.Request, td *templateData, templates ...string) error {
	// Append file path to template files
	for i, t := range templates {
		templates[i] = fmt.Sprintf("templates/%s.html", t)
	}

	// Create/get/cache template
	var t *template.Template
	var err error
	if app.Config.Env == "production" {
		// TODO: Check if file exists in template cache

		// Read template files from embedded fs
		t, err = template.New(filepath.Base(templates[0])).
			Funcs(templateFunctions).
			ParseFS(templatesFiles, templates...)

		// TODO: Store template in template cache
	} else {
		// Read template files from disk
		t, err = template.New(filepath.Base(templates[0])).
			Funcs(templateFunctions).
			ParseFiles(templates...)
	}
	if err != nil {
		// TODO
		fmt.Println(err)
		return err
	}

	// Add default template data
	if td == nil {
		td = &templateData{}
	}
	td = app.addDefaultTemplateData(td, r)

	err = t.Execute(w, td)
	if err != nil {
		log.Fatal(err)
		// app.errorLog.Println(err)
		return err
	}

	return nil
}

// Add default data to template data struct
func (app *application) addDefaultTemplateData(td *templateData, r *http.Request) *templateData {
	td.Csrf = csrf.TemplateField(r)
	return td
}
