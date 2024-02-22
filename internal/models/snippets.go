package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// A SnippetModel type to wrap the sql DB connection
type SnippetModel struct {
	DB *sql.DB
}

// Inserts new snippet data into database
func (model *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// Gets a specific snippet from database
func (model *SnippetModel) Get(id int) (*Snippet, error) {
	return &Snippet{}, nil
}

// Gets the 10 most recent snippets from the database
func (model *SnippetModel) Latest() ([]*Snippet, error) {
	return []*Snippet{}, nil
}
