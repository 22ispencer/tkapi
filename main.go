package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/WSU-Robotics-Lab/tkapi/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var env *db.Env

func main() {
	var err error
	env, err = db.GetEnv()
	if err != nil {
		log.Fatal("Unable to connect to database")
	}
	server := http.Server{
		Addr:    ":2024",
		Handler: router(),
	}
	log.Printf("Starting server at \"%s\"", server.Addr)
	server.ListenAndServe()
}

func router() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		labId, err := strconv.Atoi(r.URL.Query().Get("labId"))
		if err != nil {
			labId = 0
		}
		isActive := strings.ToLower(r.URL.Query().Get("isActive")) == "true"
		labRoleId, err := strconv.Atoi(r.URL.Query().Get("labRoleId"))
		if err != nil {
			labRoleId = 0
		}
		users, err := env.GetUsers(labId, isActive, labRoleId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error running procedure: %s", err), http.StatusInternalServerError)
		}
		if r.Header.Get("Accept") == "text/html" {
			w.Header().Add("Content-Type", "text/html")
			usersHtml := `{{range .}}<option value="{{.Id}}">{{.FirstName}} {{if .LastName.Valid}}{{.LastName.String}}{{end}}</option>{{end}}`
			tmpl := template.New("users")
			tmpl.Parse(usersHtml)

			err = tmpl.Execute(w, users)
			if err != nil {
				http.Error(w, "Error parsing users to HTML", http.StatusInternalServerError)
			}
		} else {
			w.Header().Add("Content-Type", "application/json")
			usersJson, err := json.Marshal(users)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error converting to JSON: %s", err), http.StatusInternalServerError)
			} else {
				w.Write(usersJson)
			}

		}
	})
	router.Get(`/user/{id:^\d+$}`, func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(chi.URLParam(r, "id"))
		user, err := env.GetUserById(userId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error running procedure: %s", err), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			userJson, err := json.Marshal(user)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error converting to JSON: %s", err), http.StatusInternalServerError)
			} else {
				w.Write(userJson)
			}
		}
	})
	router.Get("/labs", func(w http.ResponseWriter, r *http.Request) {
		labs, err := env.GetLabs()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error running procedure: %s", err), http.StatusInternalServerError)
		}

		if r.Header.Get("Accept") == "text/html" {
			w.Header().Add("Content-Type", "text/html")
			labsHtml := `{{range .}}<option value="{{.Id}}">{{.Name}}</option>{{end}}`
			tmpl := template.New("labs")
			tmpl.Parse(labsHtml)

			err := tmpl.Execute(w, labs)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error unable to format HTML: %s", err), http.StatusInternalServerError)
			}
		} else {
			w.Header().Add("Content-Type", "application/json")
			labsJson, err := json.Marshal(labs)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error converting to JSON: %s", err), http.StatusInternalServerError)
			} else {
				w.Write(labsJson)
			}
		}
	})
	return router
}
