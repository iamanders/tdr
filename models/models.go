package models

import (
	"database/sql"
)

var db *sql.DB

// Init model package
func InitModels(dbIn *sql.DB) {
	db = dbIn
}
