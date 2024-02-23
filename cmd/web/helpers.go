package main

import (
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

	writer.WriteHeader(status)

	if err := template.ExecuteTemplate(writer, "base", data); err != nil {
		app.serverError(writer, err)
	}
}
