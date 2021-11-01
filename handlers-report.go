package main

import (
	"log"
	"net/http"
)

// GET /report
func (app *application) GetReport(w http.ResponseWriter, r *http.Request) {
	td := templateData{Data: make(map[string]interface{})}

	if err := app.renderTemplates(w, r, &td, "report", "layout.app"); err != nil {
		log.Fatal(err)
	}
}
