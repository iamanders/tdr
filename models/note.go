package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Note model structure
type NoteModel struct {
	Id       int64
	HappenAt sql.NullTime
	Note     string
}

// To string
func (t NoteModel) String() string {
	return fmt.Sprintf("Note model #%d", t.Id)
}

// Get note for day, returns empty string or the note string
func GetNoteForDay(happenAt time.Time) (*NoteModel, error) {
	n := NoteModel{}
	dateString := happenAt.Format("2006-01-02")
	row := db.QueryRow(fmt.Sprintf("SELECT id, happen_at, note FROM notes WHERE happen_at between '%s 00:00:00' AND '%s 23:59:59' LIMIT 1", dateString, dateString))
	err := row.Scan(&n.Id, &n.HappenAt, &n.Note)
	if err == sql.ErrNoRows {
		return &n, err
	}
	return &n, nil
}

// Get notes between two dates
func GetNotes(fromDate string, toDate string) []NoteModel {
	notes := []NoteModel{}

	rows, err := db.Query(`SELECT id, happen_at, note FROM notes WHERE happen_at BETWEEN ? AND ? ORDER BY happen_at ASC`, fromDate, toDate)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		n := NoteModel{}
		err = rows.Scan(&n.Id, &n.HappenAt, &n.Note)
		if err != nil {
			log.Fatal(err)
		}
		notes = append(notes, n)
	}

	return notes
}

// Insert note to database
func InsertNoteModel(n *NoteModel) *NoteModel {
	fmt.Println(n)
	stmt, err := db.Prepare("INSERT INTO notes (happen_at, note) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var res sql.Result
	res, err = stmt.Exec(n.HappenAt.Time, n.Note)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	n.Id, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return n
}

// Store change note model
func UpdateNoteModel(n *NoteModel) *NoteModel {
	stmt, err := db.Prepare("UPDATE notes SET happen_at=?, note=? WHERE id=?")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	_, err = stmt.Exec(n.HappenAt.Time, n.Note, n.Id)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return n
}

// Delete note from database
func DeleteNoteModel(n *NoteModel) {
	if n.Id < 1 {
		return
	}

	stmt, err := db.Prepare("DELETE FROM notes WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(n.Id)
	if err != nil {
		log.Fatal(err)
	}
}
