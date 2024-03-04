package mocks

import (
	"time"

	"github.com/sparrowsl/snippetbox/internal/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "Lorem ipsum",
	Content: "Lorem ipsum dolor sit amet!",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (mock *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (mock *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (mock *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
