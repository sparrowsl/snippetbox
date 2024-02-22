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
	statement := `INSERT INTO snippets (title, content, created, expires) VALUES
            (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := model.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Get the id of the newly inserted snippet
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Gets a specific snippet from database
func (model *SnippetModel) Get(id int) (*Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets
                WHERE id = ?`

	row := model.DB.QueryRow(statement, id)
	_ = row

	return &Snippet{}, nil
}

// Gets the 10 most recent snippets from the database
func (model *SnippetModel) Latest() ([]*Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets
                  ORDER BY created DESC LIMIT 10`
	rows, err := model.DB.Query(statement)
	if err != nil {
		return nil, err
	}

	_ = rows
	return []*Snippet{}, nil
}
