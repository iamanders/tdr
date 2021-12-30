package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iamanders/tdr/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/webview/webview"
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
		Env:     "production",
		Port:    8765,
		Address: "127.0.0.1",
		CsrfKey: "00000000000000000000000000000000",
	}

	app := application{
		Config: &config,
	}

	// Flags
	flag.StringVar(&config.Env, "env", "production", "Environment, production or dev")
	flag.IntVar(&config.Port, "port", 8765, "Port to bind server to")
	flag.StringVar(&config.Address, "host", "127.0.0.1", "Host to bind server to")
	flag.StringVar(&config.CsrfKey, "csrf", "00000000000000000000000000000000", "CSRF key, 32 chars")
	flag.Parse()

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

	go func() {
		log.Printf("Starting server at http://%s:%d", config.Address, config.Port)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), app.Router))
	}()

	startOpenWindow(&config)
}

func startOpenWindow(config *config) {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("TDR")
	w.SetSize(2000, 1200, webview.HintNone)
	w.Navigate(fmt.Sprintf("http://%s:%d", config.Address, config.Port))
	// w.Navigate("http://127.0.0.1:8765")
	w.Run()
}
