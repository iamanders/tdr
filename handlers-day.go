package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iamanders/tdr/models"
)

// GET /day and /yyyy-mm-dd
func (app *application) GetDay(w http.ResponseWriter, r *http.Request) {
	td := templateData{Data: make(map[string]interface{})}

	// Default/current day or day from url?
	fromDateString := fmt.Sprintf("%s 00:00:00", time.Now().Format("2006-01-02"))
	toDateString := fmt.Sprintf("%s 23:59:59", time.Now().Format("2006-01-02"))
	urlDay := chi.URLParam(r, "date")
	if len(urlDay) > 0 {
		// Get date from url
		fromDateString = fmt.Sprintf("%s 00:00:00", urlDay)
		toDateString = fmt.Sprintf("%s 23:59:59", urlDay)
	}

	// Get time data
	td.Data["times"] = models.GetTimes(fromDateString, toDateString)
	td.Data["current"], _ = models.GetCurrentTime()

	// Get date related stuff
	currentDay, _ := time.Parse("2006-01-02 15:04:05", fromDateString)
	td.Data["day_date"] = currentDay
	td.Data["day_link_next"] = "/day/" + currentDay.AddDate(0, 0, -1).Format("2006-01-02")
	td.Data["day_link_previous"] = "/day/" + currentDay.AddDate(0, 0, 1).Format("2006-01-02")

	if err := app.renderTemplates(w, r, &td, "day", "layout.app", "partials/time-table-row"); err != nil {
		app.render500(w, r, err.Error())
		return
	}
}
