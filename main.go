package main

import (
	"log"
	"net/http"
)

// Write a home handler function which writes a byte slice as the response body
func home(writer http.ResponseWriter, req *http.Request) {
	writer.Write([]byte("Hello world from snippetbox\n"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	log.Print("Starting server on port :5000")
	err := http.ListenAndServe(":5000", mux)
	log.Fatal(err)
}
