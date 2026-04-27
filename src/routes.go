package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// Struct that stores all data that needs to be shown in a page
type Page struct {
	// Sessions
	IsLoggedIn bool

	// Theming
	//Theme string // "light" or "dark"
}

func (a *App) SlashHandler(w http.ResponseWriter, r *http.Request) {
	page := strings.TrimPrefix(r.URL.Path, "/")
	if page == "" {
		page = "home"
	}

	path := "./templates/" + page + ".html"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		a.Error(w, r, "404: Page not found! "+path)
		return
	}

	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/"+page+".html",
	)
	if err != nil {
		a.Error(w, r, fmt.Sprintf("failed to load template: %v (page=%s)", err, page))
	}

	data := Page{
		IsLoggedIn: false,
	}

	if err := tmpl.Execute(w, data); err != nil {
		a.Error(w, r, "template execution failed: ", err.Error())
	}
}

func (a *App) LoginGET(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/auth/login.html",
	)
	if err != nil {
		a.Error(w, r, fmt.Sprintf("failed to load template: %v (templates/auth/login.html)", err))
	}

	data := Page{
		IsLoggedIn: false,
	}

	if err := tmpl.Execute(w, data); err != nil {
		a.Error(w, r, "template execution failed: ", err.Error())
	}
}