package main

import (
    "fmt"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
    a = App{}
    a.Initialize("root", "", "movies-api")
    ensureTableExists()
    code := m.Run()
    clearTable()
    os.Exit(code)
}

func TestMoviesEmptyTable(t *testing.T) {
    clearTable()
    req, _ := http.NewRequest("GET", "/movies", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func TestGetNonExistentMovie(t *testing.T) {
    clearTable()
    req, _ := http.NewRequest("GET", "/movies/999", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["error"] != "Movie not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'Movie not found'. Got '%s'", m["error"])
    }
}

func TestCreateMovie(t *testing.T) {
    clearTable()
    payload := []byte(`{"titulo":"test movie","imagem":"cover.jpg","id_categoria":1,"descricao":"test movie description"}`)
    req, _ := http.NewRequest("POST", "/movies", bytes.NewBuffer(payload))
    response := executeRequest(req)
    checkResponseCode(t, http.StatusCreated, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["titulo"] != "test movie" {
        t.Errorf("Expected movie title to be 'test movie'. Got '%v'", m["titulo"])
    }

    if m["descricao"] != "test movie description" {
        t.Errorf("Expected movie description to be 'test movie description'. Got '%v'", m["descricao"])
    }

    if m["id_categoria"] != 1.0 {
        t.Errorf("Expected movie category_id to be '1'. Got '%v'", m["id_categoria"])
    }

    if m["imagem"] != "cover.jpg" {
        t.Errorf("Expected movie cover to be 'cover.jpg'. Got '%v'", m["imagem"])
    }

    // the id is compared to 1.0 because JSON unmarshaling converts numbers to
    // floats, when the target is a map[string]interface{}
    if m["id"] != 1.0 {
        t.Errorf("Expected movie ID to be '1'. Got '%v'", m["id"])
    }
}

func TestGetMovie(t *testing.T) {
    clearTable()
    addMovies(1)
    req, _ := http.NewRequest("GET", "/movies/1", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
}

func addMovies(count int) {
    if count < 1 {
        count = 1
    }

    for i := 0; i < count; i++ {
        statement := fmt.Sprintf("INSERT INTO movies(title, cover, category_id, description) VALUES('%s', '%s', %d, '%s')", ("Movie " + strconv.Itoa(i+1)), ("cover-" + strconv.Itoa(i+1) + ".jpg"), 1, ("Movie " + strconv.Itoa(i+1) + " description"))
        a.DB.Exec(statement)
    }
}

func TestUpdateMovie(t *testing.T) {
    clearTable()
    addMovies(1)
    req, _ := http.NewRequest("GET", "/movies/1", nil)
    response := executeRequest(req)

    var originalMovie map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &originalMovie)
    payload := []byte(`{"titulo":"test movie - updated name","imagem":"new-cover.jpg","id_categoria":2,"descricao":"test movie description - updated"}`)
    req, _ = http.NewRequest("PUT", "/movies/1", bytes.NewBuffer(payload))
    response = executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["id"] != originalMovie["id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalMovie["id"], m["id"])
    }

    if m["titulo"] == originalMovie["titulo"] {
        t.Errorf("Expected the title to change from '%v' to '%v'. Got '%v'", originalMovie["titulo"], m["titulo"], m["titulo"])
    }

    if m["imagem"] == originalMovie["imagem"] {
        t.Errorf("Expected the cover to change from '%v' to '%v'. Got '%v'", originalMovie["imagem"], m["imagem"], m["imagem"])
    }

    if m["id_categoria"] == originalMovie["id_categoria"] {
        t.Errorf("Expected the category_id to change from '%v' to '%v'. Got '%v'", originalMovie["id_categoria"], m["id_categoria"], m["id_categoria"])
    }

    if m["descricao"] == originalMovie["descricao"] {
        t.Errorf("Expected the description to change from '%v' to '%v'. Got '%v'", originalMovie["descricao"], m["descricao"], m["descricao"])
    }
}

func TestDeleteMovie(t *testing.T) {
    clearTable()
    addMovies(1)
    req, _ := http.NewRequest("GET", "/movies/1", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
    req, _ = http.NewRequest("DELETE", "/movies/1", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)
    req, _ = http.NewRequest("GET", "/movies/1", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func ensureTableExists() {
    if _, err := a.DB.Exec(categoriesTableCreationQuery); err != nil {
        log.Fatal(err)
    }

    if _, err := a.DB.Exec(moviesTableCreationQuery); err != nil {
        log.Fatal(err)
    }
}

func clearTable() {
    a.DB.Exec("DELETE FROM movies")
    a.DB.Exec("ALTER TABLE movies AUTO_INCREMENT = 1")
    a.DB.Exec("DELETE FROM categories")
    a.DB.Exec("ALTER TABLE categories AUTO_INCREMENT = 1")
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
