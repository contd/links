package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	var err error
	connStr := fmt.Sprintf("%s:%s@tcp(:3306)/%s", user, password, dbname)
	a.DB, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	fmt.Println("Running on http://localhost:5555/")
	log.Fatal(http.ListenAndServe(":5555", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/links", a.getLinks).Methods("GET")
	a.Router.HandleFunc("/link/{id:[0-9]+}", a.getLink).Methods("GET")
	a.Router.HandleFunc("/link", a.createLink).Methods("POST")
	a.Router.HandleFunc("/link/{id:[0-9]+}", a.updateLink).Methods("PUT")
	a.Router.HandleFunc("/link/{id:[0-9]+}", a.deleteLink).Methods("DELETE")
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
	if err := l.createLink(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
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
