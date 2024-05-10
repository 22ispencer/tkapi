package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

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
	defer env.Close()
	server := http.Server{
		Addr:    ":2024",
		Handler: router(),
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				shutdownCtxCancel()
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		fmt.Print("\r")
		log.Println("shutting down.. goodbye 👋")
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()
	log.Printf("starting server at \"%s\"", server.Addr)
	server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-serverCtx.Done()
}

func router() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
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
	router.Get(`/project/{id:^\d+$}`, func(w http.ResponseWriter, r *http.Request) {
		projectId, _ := strconv.Atoi(chi.URLParam(r, "id"))
		project, err := env.GetProjectById(projectId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error running procedure: %s", err), http.StatusInternalServerError)
		} else {
			w.Header().Add("Content-Type", "application/json")
			projectJson, err := json.Marshal(project)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error converting to JSON: %s", err), http.StatusInternalServerError)
			} else {
				w.Write(projectJson)
			}
		}
	})
	router.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
		labId, err := strconv.Atoi(r.URL.Query().Get("labId"))
		if err != nil {
			labId = 0
		}
		activeOnly := strings.ToLower(r.URL.Query().Get("activeOnly")) == "true"

		projects, err := env.GetProjects(labId, activeOnly)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error running procedure: %s", err), http.StatusInternalServerError)
		} else if r.Header.Get("Accept") == "text/html" {
			w.Header().Add("Content-Type", "text/html")
			projectsHtml := `{{range .}}<option value="{{.Id}}">{{.Name}}</option>{{end}}`
			tmpl := template.New("projects")
			tmpl.Parse(projectsHtml)

			err = tmpl.Execute(w, projects)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error unable to format HTML: %s", err), http.StatusInternalServerError)
			}
		} else {
			w.Header().Add("Content-Type", "application/json")
			projectsJson, err := json.Marshal(projects)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error converting to JSON: %s", err), http.StatusInternalServerError)
			} else {
				w.Write(projectsJson)
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
	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		labId, err := strconv.Atoi(r.URL.Query().Get("labId"))
		if err != nil {
			labId = 0
		}
		activeOnly := strings.ToLower(r.URL.Query().Get("activeOnly")) == "true"
		labRoleId, err := strconv.Atoi(r.URL.Query().Get("labRoleId"))
		if err != nil {
			labRoleId = 0
		}
		users, err := env.GetUsers(labId, activeOnly, labRoleId)
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
	return router
}
