package main

import (
    "fmt"
    "database/sql"
)

type Movie struct {
    ID int `json:"id"`
    Title string `json:"titulo"`
    Cover string `json:"imagem"`
    Category int `json:"id_categoria"`
    Description string `json:"descricao"`
}

func (m *Movie) getMovie(db *sql.DB) error {
    statement := fmt.Sprintf("SELECT title, cover, category_id, description FROM movies WHERE id=%d", m.ID)
    return db.QueryRow(statement).Scan(&m.Title, &m.Cover, &m.Category, &m.Description)
}

func (m *Movie) updateMovie(db *sql.DB) error {
    statement := fmt.Sprintf("UPDATE movies SET title='%s', cover='%s', category_id=%d, description='%s' WHERE id=%d", m.Title, m.Cover, m.Category, m.Description, m.ID)
    _, err := db.Exec(statement)
    return err
}

func (m *Movie) deleteMovie(db *sql.DB) error {
    statement := fmt.Sprintf("DELETE FROM movies WHERE id=%d", m.ID)
    _, err := db.Exec(statement)
    return err
}

func (m *Movie) createMovie(db *sql.DB) error {
    statement := fmt.Sprintf("INSERT INTO movies(title, cover, category_id, description) VALUES('%s', '%s', %d, '%s')", m.Title, m.Cover, m.Category, m.Description)

    _, err := db.Exec(statement)
    if err != nil {
        return err
    }

    err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&m.ID)
    if err != nil {
        return err
    }

    return nil
}

func getMovies(db *sql.DB, start, count int) ([]Movie, error) {
    statement := fmt.Sprintf("SELECT id, title, cover, category_id, description FROM movies LIMIT %d OFFSET %d", count, start)

    rows, err := db.Query(statement)
    if err != nil {
        return nil, err
    }

    defer rows.Close()
    movies := []Movie{}
    for rows.Next() {
        var m Movie
        if err := rows.Scan(&m.ID, &m.Title, &m.Cover, &m.Category, &m.Description); err != nil {
            return nil, err
        }
        movies = append(movies, m)
    }

    return movies, nil
}
