package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/iamanders/tdr/models"
)

// GET /report
func (app *application) GetReport(w http.ResponseWriter, r *http.Request) {
	td := templateData{Data: make(map[string]interface{})}

	// Query data
	td.Data["search_project"] = r.URL.Query().Get("project")
	td.Data["search_code"] = r.URL.Query().Get("code")
	if r.URL.Query().Get("from") != "" {
		td.Data["search_date_from"] = r.URL.Query().Get("from")
	} else {
		td.Data["search_date_from"] = time.Now().AddDate(0, 0, -6).Format("2006-01-02")
	}
	if r.URL.Query().Get("to") != "" {
		td.Data["search_date_to"] = time.Now().Format("2006-01-02")
	} else {
		td.Data["search_date_to"] = "2020-11-10"
	}

	// Search?
	td.Data["is_search"] = false
	if r.URL.Query().Get("search") == "1" {
		td.Data["is_search"] = true
		from := td.Data["search_date_from"].(string) + " 00:00:00"
		to := td.Data["search_date_to"].(string) + " 23:59:59"
		project := strings.Replace(td.Data["search_project"].(string), "*", "%", -1)
		code := strings.Replace(td.Data["search_code"].(string), "*", "%", -1)
		td.Data["search_result"] = models.GetTimesSearch(from, to, project, code)
	}

	if err := app.renderTemplates(w, r, &td, "report", "layout.app"); err != nil {
		log.Fatal(err)
	}
}
