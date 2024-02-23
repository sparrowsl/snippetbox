package models

import (
	"database/sql"
	"errors"
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
                WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := model.DB.QueryRow(statement, id)
	snippet := &Snippet{}

	// Use row.Scan() to copy the values from each field to the
	// corresponding field in the Snippet struct.
	// The arguments are pointers and don't return a value back.
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		// row.Scan() will return a sql.ErrNoRows if query returns no row.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

// Gets the 10 most recent snippets from the database that are not expired
func (model *SnippetModel) Latest() ([]*Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets
                WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := model.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	return snippets, nil
}
