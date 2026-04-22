package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/justinas/alice"
	"log/slog"
	"net/http"
	"os"
)

type App struct {
	Log *slog.Logger
	DB  *sql.DB
}

func main() {
	app := NewApp()
	mux := http.NewServeMux()
	handler := alice.New(app.RouteLogger).Then(mux)

	staticFS := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFS))

	// semi-static pages
	mux.HandleFunc("/", app.SlashHandler)

	app.Log.Info("Listening on :8000...")
	err := http.ListenAndServe(":8000", handler)
	if err != nil {
		app.Log.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func NewApp() *App {
	db, err := sql.Open("mysql", "root:H0EeLfLnO,xDEVELOPERSx4c!#%@tcp(127.0.0.1:3306)/ironite")
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		panic(err)
		os.Exit(1)
	}

	a := &App{
		Log: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		DB:  db,
	}

	return a
}
