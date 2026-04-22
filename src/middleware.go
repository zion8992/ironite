package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func (app *App) RouteLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Log.Info("new request", slog.String("url", r.URL.Path), slog.String("method", r.Method))
		next.ServeHTTP(w, r)
	})
}

func (a *App) Error(w http.ResponseWriter, r *http.Request, errs ...string) {
	msg := strings.Join(errs, ", ")

	a.Log.Error("request error",
		"path", r.URL.Path,
		"errors", msg,
	)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Error: %s", msg)
}
