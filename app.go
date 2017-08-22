package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	Title      string
	Author     string
	Categories []Category
}

type Category struct {
	Name   string
	Tag    string
	Color  string
	Border string
}

type App struct {
	Router *mux.Router
	DB     *sqlx.DB
}

func (a *App) Initialize(dbname string) {
	var err error
	a.DB, err = sqlx.Connect("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Println("Running on http://localhost:5555/")
	corsObj := handlers.AllowedOrigins([]string{"*"})
	log.Println(http.ListenAndServe(":5555", handlers.CORS(corsObj)(a.Router)))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.getIndex).Methods("GET")
	a.Router.HandleFunc("/links", a.getLinks).Methods("GET")
	a.Router.HandleFunc("/link/{id:[0-9]+}", a.getLink).Methods("GET")
	a.Router.HandleFunc("/link", a.createLink).Methods("POST")
	a.Router.HandleFunc("/link/{id:[0-9]+}", a.updateLink).Methods("PUT")
	a.Router.HandleFunc("/link/{id:[0-9]+}", a.deleteLink).Methods("DELETE")
}

func (a *App) getIndex(w http.ResponseWriter, r *http.Request) {
	cats := []Category{}
	cats = append(cats, Category{Name: "Javascript", Tag: "javascript", Color: "lightblue", Border: "blue"})
	cats = append(cats, Category{Name: "Coding", Tag: "coding", Color: "lightcoral", Border: "red"})
	cats = append(cats, Category{Name: "Tutorial", Tag: "tutorial", Color: "lightgreen", Border: "green"})
	cats = append(cats, Category{Name: "Github", Tag: "github", Color: "lightgrey", Border: "black"})
	cats = append(cats, Category{Name: "Jobs", Tag: "jobs", Color: "yellow", Border: "red"})

	page := &Page{Title: "Links Saved", Author: "Jason Kumpf", Categories: cats}
	t, _ := template.ParseFiles("links.html")
	t.Execute(w, page)
}

func (a *App) getLinks(w http.ResponseWriter, r *http.Request) {
	links, err := getLinks(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, links)
}

func (a *App) getLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid link ID")
		return
	}

	l := link{ID: id}
	if err := l.getLink(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Link not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, l)
}

func (a *App) createLink(w http.ResponseWriter, r *http.Request) {
	var l link
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&l); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	id, err := l.createLink(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	l.ID = int(id)
	respondWithJSON(w, http.StatusCreated, l)
}

func (a *App) updateLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid link ID")
		return
	}

	var l link
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&l); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()

	l.ID = id

	if err := l.updateLink(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, l)
}

func (a *App) deleteLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Link ID")
		return
	}

	l := link{ID: id}

	if err := l.deleteLink(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
