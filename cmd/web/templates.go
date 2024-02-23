package main

import "github.com/sparrowsl/snippetbox/internal/models"

// Acts as a holding structure for any dynamic data to pass in our templates
type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
