package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/justinas/alice"
	"log/slog"
	"net/http"
	"os"
	"fmt"
)

type App struct {
	Log *slog.Logger
	DB  *sql.DB
	BlockedWords string
}

func main() {
	app := NewApp()
	mux := http.NewServeMux()
	handler := alice.New(app.RouteLogger).Then(mux)

	staticFS := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFS))
	
	// auth
	mux.HandleFunc("/register", app.RegisterGET)

	// api v1
	mux.HandleFunc("/api/v1/auth/register", app.APIv1_RegisterPost)


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
		fmt.Printf("%s\n", "Failed to connect to database!")
		panic(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("%s\n", "Failed to ping to database!")
		panic(err)
	}

	a := &App{
		Log: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		DB:  db,
	}

	return a
}
