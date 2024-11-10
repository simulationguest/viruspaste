package main

import (
	"database/sql"
	"embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"

	_ "github.com/joho/godotenv/autoload"
	"leoj.de/virbin/templates"
)

//go:embed all:static
var assets embed.FS

var db *sql.DB

func routes() {
	http.Handle("/components/new-paste-form", templ.Handler(templates.NewPasteForm()))
	http.Handle("/components/new-paste-button", templ.Handler(templates.NewPasteButton()))

	http.HandleFunc("/forms/new-paste", NewPasteForm)

	http.Handle("/store/", http.StripPrefix("/store/", http.FileServer(http.Dir("./store"))))

	http.Handle("/", templ.Handler(templates.Home()))
	http.Handle("/static/", http.FileServer(http.FS(assets)))
}

func serve() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Listening on port :" + port)
	http.ListenAndServe(":"+port, nil)
}

func main() {
	defer db.Close()

	routes()

	serve()
}
