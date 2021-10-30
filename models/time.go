package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Time model structure
type TimeModel struct {
	Id       int64
	Project  string
	Code     string
	Comment  string
	StartsAt sql.NullTime
	StopsAt  sql.NullTime
	Flags    TimeModelFlag
	// Flags    int8
}

// TimeModelFlag type
type TimeModelFlag int8

// TimeModelFlag consts
const (
	TimeModelFlagNone  TimeModelFlag = 0
	TimeModelFlagRed   TimeModelFlag = 1
	TimeModelFlagGreen TimeModelFlag = 2
	TimeModelFlagBlue  TimeModelFlag = 3
)

// To string
func (t TimeModel) String() string {
	return fmt.Sprintf("Time model #%d", t.Id)
}

// Get times
func GetTimes(fromDate string, toDate string) []TimeModel {
	times := []TimeModel{}
	rows, err := db.Query(`
		SELECT id, project, code, comment, starts_at, stops_at, flags
		FROM times
		WHERE starts_at BETWEEN ? AND ?
		ORDER BY starts_at ASC`,
		fromDate, toDate,
	)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		t := TimeModel{}
		err = rows.Scan(&t.Id, &t.Project, &t.Code, &t.Comment, &t.StartsAt, &t.StopsAt, &t.Flags)
		if err != nil {
			log.Fatal(err)
		}
		times = append(times, t)
	}

	return times
}

// Get any current not completed tasks
func GetCurrentTime() (*TimeModel, error) {
	t := TimeModel{}
	row := db.QueryRow("SELECT id, project, code, comment, starts_at, stops_at, flags FROM times WHERE stops_at IS NULL ORDER BY starts_at DESC")
	err := row.Scan(&t.Id, &t.Project, &t.Code, &t.Comment, &t.StartsAt, &t.StopsAt, &t.Flags)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &t, nil
}

// Get a time by its id
func GetTimeById(id int64) (*TimeModel, error) {
	t := TimeModel{}
	row := db.QueryRow("SELECT id, project, code, comment, starts_at, stops_at, flags FROM times WHERE id = ? LIMIT 1", id)
	err := row.Scan(&t.Id, &t.Project, &t.Code, &t.Comment, &t.StartsAt, &t.StopsAt, &t.Flags)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &t, nil
}

// Insert a time to database
func InsertTimeModel(t *TimeModel) *TimeModel {
	stmt, err := db.Prepare("INSERT INTO times (project, code, comment, starts_at, stops_at, flags) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var res sql.Result
	if t.StopsAt.Valid {
		res, err = stmt.Exec(t.Project, t.Code, t.Comment, t.StartsAt.Time, t.StopsAt.Time, t.Flags)
	} else {
		res, err = stmt.Exec(t.Project, t.Code, t.Comment, t.StartsAt.Time, nil, t.Flags)
	}
	if err != nil {
		log.Fatal(err)
		return nil
	}

	t.Id, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return t
}

// Store change to database
func UpdateTimeModel(t *TimeModel) *TimeModel {

	stmt, err := db.Prepare("UPDATE times SET project=?, code=?, comment=?, starts_at=?, stops_at=?, flags=? WHERE id=?")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	_, err = stmt.Exec(t.Project, t.Code, t.Comment, t.StartsAt.Time, t.StopsAt, t.Flags, t.Id)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return t
}

// Delete time
func DeleteTimeModel(t *TimeModel) {
	if t.Id < 1 {
		return
	}

	stmt, err := db.Prepare("DELETE FROM times WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(t.Id)
	if err != nil {
		log.Fatal(err)
	}
}

// Return time difference between two time models
func (t *TimeModel) TimeDuration() time.Duration {
	if t.StartsAt.Valid && t.StopsAt.Valid {
		return t.StopsAt.Time.Sub(t.StartsAt.Time)
	} else {
		return time.Now().UTC().Sub(t.StartsAt.Time)
	}
}

// Return total time difference between multiple time models
func TimeDurationMultiple(times []TimeModel, includeRunning bool) time.Duration {
	var td time.Duration
	for _, t := range times {
		if t.StartsAt.Valid {
			if includeRunning || (!includeRunning && t.StopsAt.Valid) {
				td += t.TimeDuration()
			}
		}
	}
	return td
}
