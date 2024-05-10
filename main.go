package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	server := http.Server{
		Addr:    "2024",
		Handler: router(),
	}
	server.ListenAndServe()
}

func router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	r.Get("users", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Accept"), "application/json") {
			w.Header().Add("Content-Type", "application/json")
		}
	})
	return r
}
