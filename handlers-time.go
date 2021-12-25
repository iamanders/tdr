package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iamanders/tdr/models"
)

// GET Time edit
func (app *application) GetTimeEdit(w http.ResponseWriter, r *http.Request) {

	timeId := chi.URLParam(r, "timeId")
	timeIdInt, err := strconv.Atoi(timeId)
	if err != nil {
		log.Fatal(err)
		// TODO: return something
	}

	t, err := models.GetTimeById(int64(timeIdInt))
	if err != nil {
		log.Fatal(err)
		// TODO: return something
	}

	td := templateData{Data: make(map[string]interface{})}
	td.Data["time"] = t
	if err := app.renderTemplates(w, r, &td, "edit", "layout.app"); err != nil {
		log.Fatal(err)
		// app.errorLog.Println(err)
	}
}

// POST Time store
func (app *application) PostTimeStore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t := models.TimeModel{
		Project:  r.FormValue("project"),
		Code:     r.FormValue("code"),
		Comment:  r.FormValue("comment"),
		StartsAt: sql.NullTime{Time: time.Now(), Valid: true},
		Flags:    models.TimeModelFlagNone,
	}

	// Optional - Starts at
	if len(r.FormValue("starts_at")) > 0 {
		d, err := time.ParseInLocation("2006-01-02 15:04", r.FormValue("starts_at"), time.Now().Location())
		if err == nil {
			t.StartsAt = sql.NullTime{Time: d, Valid: true}
		}
	}

	// Optional - Stops at
	if len(r.FormValue("stops_at")) > 0 {
		d, err := time.ParseInLocation("2006-01-02 15:04", r.FormValue("stops_at"), time.Now().Location())
		if err == nil {
			t.StopsAt = sql.NullTime{Time: d, Valid: true}
		}
	}

	// Optional - No gap?
	if r.FormValue("no_gap") == "1" {
		lastTime, err := models.GetLastTimeOfDay(time.Now().Format("2006-01-02"))
		if err == nil && lastTime.Id > 0 {
			t.StartsAt = lastTime.StopsAt
		}
	}

	// TODO: Validate
	// TEMP: Before validation is done
	defaultEmptyString(&t.Project, "Default project")

	// Insert and redirect
	models.InsertTimeModel(&t)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

// POST Time update
func (app *application) PostTimeStop(w http.ResponseWriter, r *http.Request) {
	// Load time
	timeId := chi.URLParam(r, "timeId")
	timeIdInt, err := strconv.Atoi(timeId)
	if err != nil {
		log.Fatal(err)
		// TODO: return something
	}
	t, err := models.GetTimeById(int64(timeIdInt))
	if err != nil {
		log.Fatal(err)
		// TODO: return something
	}

	// Update
	t.StopsAt = sql.NullTime{Time: time.Now(), Valid: true}
	models.UpdateTimeModel(t)

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

// Update time
func (app *application) PostTimeUpdate(w http.ResponseWriter, r *http.Request) {
	// Load time
	timeId := chi.URLParam(r, "timeId")
	timeIdInt, err := strconv.Atoi(timeId)
	if err != nil {
		log.Fatal(err)
		// TODO: return something
	}
	t, err := models.GetTimeById(int64(timeIdInt))
	if err != nil {
		log.Fatal(err)
		// TODO: return something
	}

	// Get data from form
	r.ParseForm()
	t.Project = r.FormValue("project")
	t.Code = r.FormValue("code")
	t.Comment = r.FormValue("comment")
	t.Flags = models.TimeModelFlag(intOrZero(r.FormValue("flags")))

	dateLayout := "2006-01-02 15:04"

	startsAtString := fmt.Sprintf("%s %s", r.FormValue("starts_at_date"), r.FormValue("starts_at_time"))
	startsAtDate, err := time.ParseInLocation(dateLayout, startsAtString, time.Now().Location())
	if err != nil {
		t.StartsAt = sql.NullTime{}
	} else {
		t.StartsAt.Time = startsAtDate
		t.StartsAt.Valid = true
	}

	stopsAtString := fmt.Sprintf("%s %s", r.FormValue("stops_at_date"), r.FormValue("stops_at_time"))
	stopsAtDate, err := time.ParseInLocation(dateLayout, stopsAtString, time.Now().Location())
	if err != nil {
		t.StopsAt = sql.NullTime{Valid: false}
	} else {
		t.StopsAt.Valid = true
		t.StopsAt.Time = stopsAtDate
	}

	// Update
	models.UpdateTimeModel(t)

	// Redirect
	dateString := t.StartsAt.Time.Format("2006-01-02")
	http.Redirect(w, r, "/day/"+dateString, http.StatusFound)
	// http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

// DELETE time
func (app *application) PostTimeDelete(w http.ResponseWriter, r *http.Request) {
	// Load time
	timeId := chi.URLParam(r, "timeId")
	timeIdInt, err := strconv.Atoi(timeId)
	if err != nil {
		log.Fatal(err)
		// TODO: return something
	}
	t, err := models.GetTimeById(int64(timeIdInt))
	if err != nil {
		log.Fatal(err)
		// TODO: return something
	}

	// Delete
	dateString := t.StartsAt.Time.Format("2006-01-02")
	models.DeleteTimeModel(t)
	http.Redirect(w, r, "/day/"+dateString, http.StatusFound)
}
