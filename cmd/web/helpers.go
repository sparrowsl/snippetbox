package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(writer http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
}

func (app *application) notFound(writer http.ResponseWriter) {
	app.clientError(writer, http.StatusNotFound)
}

func (app *application) render(writer http.ResponseWriter, status int, page string, data *TemplateData) {
	template, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s doesn't exist!", page)
		app.serverError(writer, err)
		return
	}

	// Init a new temporary buffer to write to
	buffer := new(bytes.Buffer)

	// Write data to buffer. If error, return server error.
	if err := template.ExecuteTemplate(buffer, "base", data); err != nil {
		app.serverError(writer, err)
		return
	}

	// If no error, then set header and move data from buffer to the writer
	writer.WriteHeader(status)
	buffer.WriteTo(writer)
}

func (app *application) Authenticate(request *http.Request) bool {
	return app.sessionManager.Exists(request.Context(), "authenticatedUserID")
}
