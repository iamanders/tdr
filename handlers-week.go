package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iamanders/tdr/models"
	"github.com/snabb/isoweek"
)

// summaryData represents project and times
// Used on week page
type summaryData struct {
	Project  string
	Duration time.Duration
	Codes    map[string]time.Duration
}

// summaryFindProject finds a project in a []summaryData array
// Returns index of project. Returns -1 if not found.
func summaryFindProject(rows *[]summaryData, project string) int {
	for i, row := range *rows {
		if row.Project == project {
			return i
		}
	}
	return -1
}

// GET week index
func (app *application) GetWeek(w http.ResponseWriter, r *http.Request) {
	td := templateData{Data: make(map[string]interface{})}

	// Default date
	defaultYear, defaultWeek := time.Now().ISOWeek()
	startDate := isoweek.StartTime(defaultYear, defaultWeek, time.UTC)

	// Query parameter date?
	urlYear := chi.URLParam(r, "year")
	urlWeek := chi.URLParam(r, "week")
	if len(urlYear) > 0 && len(urlWeek) > 0 {
		week, _ := strconv.Atoi(urlWeek)
		year, _ := strconv.Atoi(urlYear)
		startDate = isoweek.StartTime(year, week, time.UTC)
	}
	stopDate := startDate.AddDate(0, 0, 7).Add(time.Second * -1) // Add 7 days and remove one second

	// Template data
	td.Data["start_date"] = startDate
	td.Data["stop_date"] = stopDate
	td.Data["year"], td.Data["week"] = startDate.ISOWeek()
	td.Data["next_year"], td.Data["next_week"] = startDate.AddDate(0, 0, 7).ISOWeek()
	td.Data["previous_year"], td.Data["previous_week"] = startDate.AddDate(0, 0, -7).ISOWeek()

	// Times
	td.Data["times"] = models.GetTimes(startDate.Format("2006-01-02 15:04:05"), stopDate.Format("2006-01-02 15:04:05"))

	//  Summary per project and code
	summary := make([]summaryData, 0)
	// Make sure all projects exists in summary array
	for _, t := range td.Data["times"].([]models.TimeModel) {
		rowIndex := summaryFindProject(&summary, t.Project)
		if rowIndex < 0 {
			row := &summaryData{Project: t.Project, Duration: 0, Codes: make(map[string]time.Duration)}
			summary = append(summary, *row)
		}
	}
	// Append summary data
	for _, t := range td.Data["times"].([]models.TimeModel) {
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

	// Render
	if err := app.renderTemplates(w, r, &td, "week", "layout.app", "partials/time-table-row"); err != nil {
		log.Fatal(err)
		// app.errorLog.Println(err)
	}
}
