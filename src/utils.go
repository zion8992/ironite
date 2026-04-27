package main

import(
	"fmt"
	"strings"
	"net/http"
)

func (a *App) Error(w http.ResponseWriter, r *http.Request, errs ...string) {
	msg := strings.Join(errs, "")

	a.Log.Error("request error",
		"path", r.URL.Path,
		"errors", msg,
	)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Error: %s", msg)
}
