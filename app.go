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
    "os"
)

type App struct {
    Router *mux.Router
    DB *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
    connectionString := os.Getenv("DATABASE_URL")

    if (connectionString == "") {
        connectionString = fmt.Sprintf("%s:%s@/%s", user, password, dbname)
    }

    var err error
    a.DB, err = sql.Open("mysql", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    const categoriesTableCreationQuery = `
    CREATE TABLE IF NOT EXISTS categories
    (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(50) NOT NULL
    )`

    const moviesTableCreationQuery = `
    CREATE TABLE IF NOT EXISTS movies
    (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(120) NOT NULL,
        cover VARCHAR(255),
        category_id INT NOT NULL,
        description text,
        FOREIGN KEY (category_id) REFERENCES Categories(id)
    )`

    if _, err := a.DB.Exec(categoriesTableCreationQuery); err != nil {
        log.Fatal(err)
    }

    if _, err := a.DB.Exec(moviesTableCreationQuery); err != nil {
        log.Fatal(err)
    }

    a.Router = mux.NewRouter()
    a.initializeRoutes()
}

func (a *App) Run(addr string) {
    log.Println(addr)
    log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/movies", a.getMovies).Methods("GET")
    a.Router.HandleFunc("/movies", a.createMovie).Methods("POST")
    a.Router.HandleFunc("/movies/{id:[0-9]+}", a.getMovie).Methods("GET")
    a.Router.HandleFunc("/movies/{id:[0-9]+}", a.updateMovie).Methods("PUT")
    a.Router.HandleFunc("/movies/{id:[0-9]+}", a.deleteMovie).Methods("DELETE")
}

func (a *App) getMovie(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
        return
    }
    m := Movie{ID: id}
    if err := m.getMovie(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Movie not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }
    respondWithJSON(w, http.StatusOK, m)
}

func (a *App) getMovies(w http.ResponseWriter, r *http.Request) {
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))
    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }
    movies, err := getMovies(a.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusOK, movies)
}

func (a *App) createMovie(w http.ResponseWriter, r *http.Request) {
    var m Movie
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&m); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
    if err := m.createMovie(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusCreated, m)
}

func (a *App) updateMovie(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid movie ID")
        return
    }
    var m Movie
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&m); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
        return
    }
    defer r.Body.Close()
    m.ID = id
    if err := m.updateMovie(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusOK, m)
}

func (a *App) deleteMovie(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid User ID")
        return
    }
    m := Movie{ID: id}
    if err := m.deleteMovie(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}
