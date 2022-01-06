package main

import (
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
		td.Data["search_date_from"] = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	if r.URL.Query().Get("to") != "" {
		td.Data["search_date_to"] = r.URL.Query().Get("to")
	} else {
		td.Data["search_date_to"] = time.Now().Format("2006-01-02")
	}

	// Search?
	td.Data["is_search"] = false
	if r.URL.Query().Get("search") == "1" {
		// Set vars
		td.Data["is_search"] = true
		from := td.Data["search_date_from"].(string) + " 00:00:00"
		to := td.Data["search_date_to"].(string) + " 23:59:59"
		project := strings.Replace(td.Data["search_project"].(string), "*", "%", -1)
		code := strings.Replace(td.Data["search_code"].(string), "*", "%", -1)
		td.Data["search_result"] = models.GetTimesSearch(from, to, project, code)
		td.Data["notes"] = models.GetNotes(from, to)

		//  Summary per project and code
		summary := make([]summaryData, 0)
		// Make sure all projects exists in summary array
		for _, t := range td.Data["search_result"].([]models.TimeModel) {
			rowIndex := summaryFindProject(&summary, t.Project)
			if rowIndex < 0 {
				row := &summaryData{Project: t.Project, Duration: 0, Codes: make(map[string]time.Duration)}
				summary = append(summary, *row)
			}
		}
		// Append summary data
		for _, t := range td.Data["search_result"].([]models.TimeModel) {
			rowIndex := summaryFindProject(&summary, t.Project)
			summary[rowIndex].Duration += t.TimeDuration()

			// Get or create code
			if _, ok := summary[rowIndex].Codes[t.Code]; !ok {
				// Init the code row
				summary[rowIndex].Codes[t.Code] = 0
			}
			summary[rowIndex].Codes[t.Code] += t.TimeDuration()
		}
		td.Data["summary"] = summary
	}

	if err := app.renderTemplates(w, r, &td, "report", "layout.app", "partials/time-table-row"); err != nil {
		app.render500(w, r, err.Error())
		return
	}
}
