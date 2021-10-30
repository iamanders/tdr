package main

import (
	"net/http"
)

// GET /
func (app *application) GetHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/day", http.StatusFound)
}
