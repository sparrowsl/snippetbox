package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Write a home handler function which writes a byte slice as the response body
func home(writer http.ResponseWriter, request *http.Request) {
	// If request is not home then show 404 page
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	// read template file into a template set with template.ParseFiles()
	temp, err := template.ParseFiles("./ui/html/pages/home.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Use Execute() to write template as content for the response body
	if err := temp.Execute(writer, nil); err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// A handler to create new snippet
func createSnippet(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		// Set supported format for URL - Allowed methods; POST
		writer.Header().Add("Allow", http.MethodPost)
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	writer.Write([]byte("Create new snippet!"))
}

// A handler to view a specific snippet
func viewSnippet(writer http.ResponseWriter, request *http.Request) {
	// Get the query params 'id' from the request
	queryId := request.URL.Query().Get("id")

	// Check if 'id' is valid - by converting from string to int
	// or if id is not less than 0
	id, err := strconv.Atoi(queryId)
	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}

	fmt.Fprintf(writer, "Displaying snippet with the id of %d", id)
}
