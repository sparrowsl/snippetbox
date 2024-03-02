package main

import (
	"html/template"
	"path/filepath"

	"github.com/sparrowsl/snippetbox/internal/models"
)

// Acts as a holding structure for any dynamic data to pass in our templates
type TemplateData struct {
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Errors          map[string]string
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
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

		// Parse the base template file
		temp, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// call ParseGlob() on the base template to add partials
		temp, err = temp.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// call ParseFiles() to add the extra page on the template set
		temp, err = temp.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template to the cache/map as normal
		cache[name] = temp
	}

	return cache, nil
}
