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

// Page is the structure for saved urls or pages.
type Page struct {
	Title      string
	Author     string
	Categories []Category
}

// Category is the main categorie for saved urls.
type Category struct {
	Name  string
	Tag   string
	Color string
	Font  string
}

// App is the main app in a struct for unit testability.
type App struct {
	Router *mux.Router
	DB     *sqlx.DB
}

// Initialize the app and the database connection. Better than init()
func (a *App) Initialize(dbname string) {
	var err error
	a.DB, err = sqlx.Connect("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run is how the server is started, so other parts can be unit tested.
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
	cats = append(cats, Category{Name: "Shell", Tag: "shell", Color: "#001f3f", Font: "white"})
	cats = append(cats, Category{Name: "Awesome", Tag: "awesome", Color: "#0074D9", Font: "white"})
	cats = append(cats, Category{Name: "CSS", Tag: "css", Color: "#7FDBFF", Font: "#001f3f"})
	cats = append(cats, Category{Name: "DevOps", Tag: "devops", Color: "#39CCCC", Font: "#001f3f"})
	cats = append(cats, Category{Name: "Data", Tag: "data", Color: "#3D9970", Font: "#001f3f"})
	cats = append(cats, Category{Name: "Github", Tag: "github", Color: "#2ECC40", Font: "#001f3f"})
	cats = append(cats, Category{Name: "Go", Tag: "go", Color: "#01FF70", Font: "#001f3f"})
	cats = append(cats, Category{Name: "Hack", Tag: "hack", Color: "#FFDC00", Font: "#001f3f"})
	cats = append(cats, Category{Name: "Javascript", Tag: "javascript", Color: "#FF851B", Font: "#001f3f"})
	cats = append(cats, Category{Name: "Work", Tag: "shell", Color: "#B10DC9", Font: "white"})
	cats = append(cats, Category{Name: "Sites", Tag: "sites", Color: "#85144b", Font: "white"})
	cats = append(cats, Category{Name: "Tools", Tag: "tools", Color: "#DDDDDD", Font: "#001f3f"})

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
