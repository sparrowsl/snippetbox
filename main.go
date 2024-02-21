package main

import (
	"log"
	"net/http"
)

// Write a home handler function which writes a byte slice as the response body
func home(writer http.ResponseWriter, req *http.Request) {
	// If request is not home then show 404 page
	if req.URL.Path != "/" {
		http.NotFound(writer, req)
		return
	}

	writer.Write([]byte("Hello world from snippetbox"))
}

// A handler to create new snippet
func createSnippet(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		// Set supported format for URL - Allowed methods; POST
		writer.Header().Add("Allow", "POST")
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	writer.Write([]byte("Create new snippet!"))
}

// A handler to view a specific snippet
func viewSnippet(writer http.ResponseWriter, req *http.Request) {
	writer.Write([]byte("Display a specific snippet"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", viewSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Print("Starting server on port :5000")
	err := http.ListenAndServe(":5000", mux)
	log.Fatal(err)
}
