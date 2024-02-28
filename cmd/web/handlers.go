package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
	"github.com/sparrowsl/snippetbox/internal/models"
)

// Write a home handler function which writes a byte slice as the response body
func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, http.StatusOK, "home.html", &TemplateData{
		Snippets: snippets,
	})
}

// A handler to display form to create snippets
func (app *application) createSnippet(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, http.StatusOK, "create.html", nil)
}

// A handler to create new snippet
func (app *application) createSnippetPost(writer http.ResponseWriter, request *http.Request) {
	// parse form data
	if err := request.ParseForm(); err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	title := request.PostForm.Get("title")
	content := request.PostForm.Get("content")
	expires, err := strconv.Atoi(request.PostForm.Get("expires"))
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	// Do validation checks for incoming data
	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must equal to 1, 7, or 365"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprint(writer, fieldErrors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// A handler to view a specific snippet
func (app *application) viewSnippet(writer http.ResponseWriter, request *http.Request) {
	queryId := chi.URLParam(request, "id")

	id, err := strconv.Atoi(queryId)
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}

	app.render(writer, http.StatusOK, "view.html", &TemplateData{
		Snippet: snippet,
	})
}
