package main

import (
    "database/sql"
    "errors"
)

type Category struct {
    ID int `json:"id"`
    Title string `json:"titulo"`
}

func (m *Category) getCategory(db *sql.DB) error {
    return errors.New("Not implemented yet");
}

func (m *Category) updateCategory(db *sql.DB) error {
    return errors.New("Not implemented yet");
}

func (m *Category) deleteCategory(db *sql.DB) error {
    return errors.New("Not implemented yet");
}

func (m *Category) createCategory(db *sql.DB) error {
    return errors.New("Not implemented yet");
}

func getCategories(db *sql.DB, start, count int) ([]Category, error) {
    return nil, errors.New("Not implemented yet");
}
