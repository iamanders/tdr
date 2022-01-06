package main

import (
	"os"
)

// Create SQLite database file
func instantiateDatabaseFile(filename string) error {
	// Create file and then close it
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

// Create database file structure if it is not already created
func (app *application) setupDatabaseStructure() error {
	// Setup times table
	_, err := app.DB.Exec("CREATE TABLE IF NOT EXISTS `times` (`id` integer, `project` varchar, `code` varchar, `comment` varchar, `starts_at` datetime, `stops_at` datetime, `flags` int, PRIMARY KEY (id));")
	if err != nil {
		return err
	}

	// Setup notes table
	_, err = app.DB.Exec("CREATE TABLE IF NOT EXISTS `notes` (`id` integer, `happen_at` datetime, `note` text, PRIMARY KEY (id));")
	if err != nil {
		return err
	}

	return nil
}
