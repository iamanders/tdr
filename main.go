package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iamanders/tdr/models"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed templates
var templatesFiles embed.FS

//go:embed static
var staticFiles embed.FS

// Holds application config
type config struct {
	Env     string
	Port    int
	Address string
	CsrfKey string
}

// Holds application variables
type application struct {
	Config    *config
	Router    *chi.Mux
	DB        *sql.DB
	Templates map[string]*template.Template
}

func main() {

	// Config and application variable
	config := config{
		Env: "dev",
		// Env:     "production",
		Port:    8765,
		Address: "127.0.0.1",
		CsrfKey: "00000000000000000000000000000000",
	}

	app := application{
		Config: &config,
	}

	// Environment variables
	if environment := os.Getenv("ENV"); len(environment) > 0 {
		config.Env = environment
	}
	if port := os.Getenv("PORT"); len(port) > 0 {
		config.Port, _ = strconv.Atoi(port)
	}
	if address := os.Getenv("HOST"); len(address) > 0 {
		config.Address = address
	}
	if csrf := os.Getenv("CSRF"); len(csrf) > 0 {
		config.CsrfKey = csrf
	}

	// DB
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dbFilename := fmt.Sprintf("%s/.tdr.db", homedir)

	err = instantiateDatabaseFile(dbFilename)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	app.DB = db

	err = app.setupDatabaseStructure()
	if err != nil {
		log.Fatal(err)
	}

	// Models
	models.InitModels(app.DB)

	// Template cache
	app.Templates = make(map[string]*template.Template)

	// Router
	app.Router = chi.NewRouter()
	app.SetupRoutes()

	// Start server
	log.Printf("Starting server at http://%s:%d", config.Address, config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), app.Router))
}
