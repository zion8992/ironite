package main

/** API V1
Documentation at /docs/api/v1/index.md
**/

import(
	"net/http"
	"fmt"
	"html/template"
)

func (a *App) APIv1_RegisterPost(w http.ResponseWriter, r *http.Request) {
	if a.WantsJson(r)

	r.ParseForm()
	ok, err := a.ValidateFormFields(r.FormValue("username"), r.FormValue("password"), r.FormValue("email"))

	var failMessage string
	if !ok {
		failMessage = "Form field validation failed: "+err.Error()
	}

	tmpl, err := template.ParseFiles(
		"./templates/base.html",
		"./templates/auth/register_ok.html",
	)
	if err != nil {
		a.Error(w, r, fmt.Sprintf("failed to load template: %v (templates/auth/login.html)", err))
	}

	data := Page{
		IsLoggedIn: false,
		FailMessage: failMessage,
	}

	if err := tmpl.Execute(w, data); err != nil {
		a.Error(w, r, "template execution failed: ", err.Error())
	}
}