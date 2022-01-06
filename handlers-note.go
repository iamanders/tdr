package main

import (
	"database/sql"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/iamanders/tdr/models"
)

// POST Note upsert
func (app *application) PostNoteUpsert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Parse date
	d, err := time.ParseInLocation("2006-01-02", r.FormValue("date"), time.Now().Location())
	if err != nil {
		app.render404(w, r, err.Error())
		return
	}

	// Note
	note := r.FormValue("note")

	// Get note from database
	n, err := models.GetNoteForDay(d)

	// Delete note?
	if utf8.RuneCountInString(note) < 1 && n.Id > 0 {
		// Delete note
		models.DeleteNoteModel(n)
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		return
	}

	// Upsert
	if err != nil {
		// Insert
		newNote := models.NoteModel{
			HappenAt: sql.NullTime{Time: d, Valid: true},
			Note:     note,
		}
		models.InsertNoteModel(&newNote)
	} else {
		// Update
		n.Note = note
		models.UpdateNoteModel(n)
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}
