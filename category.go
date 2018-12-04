package main

import (
    "fmt"
    "database/sql"
)

type Category struct {
    ID int `json:"id"`
    Title string `json:"titulo"`
}

type Catalog struct {
    ID int `json:"id"`
    Title string `json:"titulo"`
    Movies []Movie `json:"filmes"`
}

func (c *Category) getCategory(db *sql.DB) error {
    statement := fmt.Sprintf("SELECT title FROM categories WHERE id=%d", c.ID)
    return db.QueryRow(statement).Scan(&c.Title)
}

func (c *Category) updateCategory(db *sql.DB) error {
    statement := fmt.Sprintf("UPDATE categories SET title='%s' WHERE id=%d", c.Title, c.ID)
    _, err := db.Exec(statement)
    return err
}

func (c *Category) deleteCategory(db *sql.DB) error {
    statement := fmt.Sprintf("DELETE FROM categories WHERE id=%d", c.ID)
    _, err := db.Exec(statement)
    return err
}

func (c *Category) createCategory(db *sql.DB) error {
    statement := fmt.Sprintf("INSERT INTO categories(title) VALUES('%s')", c.Title)

    _, err := db.Exec(statement)
    if err != nil {
        return err
    }

    err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&c.ID)
    if err != nil {
        return err
    }

    return nil
}

func getCategories(db *sql.DB) ([]Category, error) {
    statement := fmt.Sprintf("SELECT id, title FROM categories")

    rows, err := db.Query(statement)
    if err != nil {
        return nil, err
    }

    defer rows.Close()
    categories := []Category{}
    for rows.Next() {
        var c Category
        if err := rows.Scan(&c.ID, &c.Title); err != nil {
            return nil, err
        }
        categories = append(categories, c)
    }

    return categories, nil
}

func getCategoriesWithMovies(db *sql.DB) ([]Catalog, error) {
    categories, err := getCategories(db)
    catalogs := []Catalog{}

    if err != nil {
        return nil, err
    }

    for _, element := range categories {
        movies, err := getMoviesByCategoryId(db, element.ID)

        if err != nil {
            return nil, err
        }

        catalogs = append(catalogs, Catalog{element.ID, element.Title, movies})
    }

    return catalogs, nil
}
