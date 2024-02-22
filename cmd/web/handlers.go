package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Write a home handler function which writes a byte slice as the response body
func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	// If request is not home then show 404 page
	if request.URL.Path != "/" {
		app.notFound(writer)
		return
	}

	// List of templates files to parse
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	// read template files into a template set with template.ParseFiles()
	temp, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	// Use ExecuteTemplate() to write the content of "base" as the response body
	// Other template (html) files will inherit from the "base" template
	if err := temp.ExecuteTemplate(writer, "base", nil); err != nil {
		app.serverError(writer, err)
	}
}

// A handler to create new snippet
func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		// Set supported format for URL - Allowed methods; POST
		writer.Header().Add("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

// A handler to view a specific snippet
func (app *application) viewSnippet(writer http.ResponseWriter, request *http.Request) {
	// Get the query params 'id' from the request
	queryId := request.URL.Query().Get("id")

	// Check if 'id' is valid - by converting from string to int
	// or if id is not less than 0
	id, err := strconv.Atoi(queryId)
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	fmt.Fprintf(writer, "Displaying snippet with the id of %d", id)
}
