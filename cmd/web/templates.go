package main

import (
	"html/template"
	"path/filepath"

	"github.com/sparrowsl/snippetbox/internal/models"
)

// Acts as a holding structure for any dynamic data to pass in our templates
type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Match all files in the filepath
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Return the name of the file from filepath
		name := filepath.Base(page)

		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			page,
		}

		temp, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = temp
	}

	return cache, nil
}
