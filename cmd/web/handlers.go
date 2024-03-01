package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sparrowsl/snippetbox/internal/models"
	"github.com/sparrowsl/snippetbox/internal/validator"
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
	app.render(writer, http.StatusOK, "create.html", &TemplateData{
		Errors: map[string]string{},
	})
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
	val := validator.Validator{}

	val.CheckField(validator.NotBlank(title), "title", "This field cannot be blank")
	val.CheckField(validator.MaxChars(title, 100), "title", "This field cannot be more than 100 characters long")
	val.CheckField(validator.NotBlank(content), "content", "This field cannot be blank")
	val.CheckField(validator.PermittedInt(expires, 7, 1, 365), "expires", "This field must be equal 1, 7, or 365")

	if !val.Valid() {
		app.render(writer, http.StatusUnprocessableEntity, "create.html", &TemplateData{
			Errors: val.FieldErrors,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.sessionManager.Put(request.Context(), "flash", "Snippet successfully created!!")

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

	flash := app.sessionManager.PopString(request.Context(), "flash")

	app.render(writer, http.StatusOK, "view.html", &TemplateData{
		Snippet: snippet,
		Flash:   flash,
	})
}

func (app *application) userSignUp(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, http.StatusOK, "signup.html", &TemplateData{
		Errors: map[string]string{},
	})
}

func (app *application) userSignUpPost(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	name := request.PostForm.Get("name")
	email := request.PostForm.Get("email")
	password := request.PostForm.Get("password")

	// validate form data
	val := validator.Validator{}

	val.CheckField(validator.NotBlank(name), "name", "This field must be less than 30 characters")
	val.CheckField(validator.NotBlank(email), "email", "This field must not be empty")
	val.CheckField(validator.Matches(email, validator.EmailRegex), "email", "This field must be a valid email address")
	val.CheckField(validator.NotBlank(password), "password", "This field must not be empty")
	val.CheckField(validator.MinChars(password, 8), "password", "This field must be at least 8 characters long")

	if !val.Valid() {
		app.render(writer, http.StatusUnprocessableEntity, "signup.html", &TemplateData{
			Errors: val.FieldErrors,
		})
		return
	}

	err := app.users.Insert(name, email, password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			val.AddFieldError("email", "Email address is already in use")
			app.render(writer, http.StatusUnprocessableEntity, "signup.html", &TemplateData{
				Errors: val.FieldErrors,
			})
		} else {
			app.serverError(writer, err)
		}
	}

	app.sessionManager.Put(request.Context(), "flash", "Your signup was successfully, please log in!")
	http.Redirect(writer, request, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(writer http.ResponseWriter, request *http.Request) {
	app.render(writer, http.StatusOK, "login.html", &TemplateData{
		Errors: map[string]string{},
	})
}

func (app *application) userLoginPost(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	name := request.PostForm.Get("name")
	password := request.PostForm.Get("password")

	// validate form data
	val := validator.Validator{}

	val.CheckField(validator.MaxChars(name, 30), "name", "This field must be less than 30 characters")
	val.CheckField(validator.NotBlank(name), "name", "This field must not be empty")
	val.CheckField(validator.NotBlank(password), "password", "This field must not be empty")

	if !val.Valid() {
		app.render(writer, http.StatusUnprocessableEntity, "login.html", &TemplateData{
			Errors: val.FieldErrors,
		})
		return
	}
}

func (app *application) userLogout(writer http.ResponseWriter, request *http.Request) {}
