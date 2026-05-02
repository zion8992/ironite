package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"database/sql"
	"time"
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
		IsLoggedIn: a.HasSessionToken(r),
	}

	if err := tmpl.Execute(w, data); err != nil {
		a.Error(w, r, "template execution failed: ", err.Error())
	}
}

func (a *App) RegisterGET(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/auth/register.html",
	)
	if err != nil {
		a.Error(w, r, fmt.Sprintf("failed to load template: %v (templates/auth/login.html)", err))
	}

	data := Page{
		IsLoggedIn: a.HasSessionToken(r),
	}

	if err := tmpl.Execute(w, data); err != nil {
		a.Error(w, r, "template execution failed: ", err.Error())
	}
}

func (a *App) RegisterPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ok, err := a.ValidateFormFields(r.FormValue("username"), r.FormValue("password"), r.FormValue("email"))

	if !ok {
		a.Error(w, r, "Form field validation failed: "+err.Error())
		return
	}

	_, err = a.GetUserIDByUsername(r.FormValue("username"))
	if err != nil && err != sql.ErrNoRows {
		a.Error(w, r, "Failed to check user: "+err.Error())
		return
	}

	if err == nil {
		a.Error(w, r, "User already exists")
		return
	}

	var passHash string
	passHash, err = a.HashPassword(r.FormValue("password"))
	if err != nil {
		a.Error(w, r, "Failed to hash password! ", err.Error())
		return
	}

	u := User {
		Username: r.FormValue("username"),
		Password: passHash,
		Email: r.FormValue("email"),
	}

	_, err = a.CreateUser(u)
	if err != nil {
		a.Error(w, r, "Failed to create user! ", err.Error())
		return
	}

	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/auth/register_ok.html",
	)
	if err != nil {
		a.Error(w, r, fmt.Sprintf("failed to load template: %v (templates/auth/login.html)", err))
		return
	}

	data := Page{
		IsLoggedIn: a.HasSessionToken(r),
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
		IsLoggedIn: a.HasSessionToken(r),
	}

	if err := tmpl.Execute(w, data); err != nil {
		a.Error(w, r, "template execution failed: ", err.Error())
	}
}

func (a *App) LoginPOST(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id, err := a.GetUserIDByUsername(r.FormValue("username"))
	if err == sql.ErrNoRows {
		a.Error(w, r, "Invalid credentials")
		return
	}
	if err != nil {
		a.Error(w, r, "Failed to check user: "+err.Error())
		return
	}

	var ok bool
	ok, err = a.CheckPasswordDB(id, r.FormValue("password"))
	if err != nil {
		a.Error(w, r, "Failed to check if your password is correct: " + err.Error())
	}

	if !ok {
		a.Error(w, r, "Invalid credentials")
		return
	}

	var token string
	token, err = a.GenerateSessionToken()
	if err != nil {
		a.Error(w, r, "Failed to create a session token to log you in: ", err.Error())
		return
	}

	expires := time.Now().Add(a.DefaultExpiry)
	err = a.SetSessionToken(id, token, expires)

	c := http.Cookie {
		Name: "session_token",
		Value: token,
		Expires: expires,
	}
	
	http.SetCookie(w, &c)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}