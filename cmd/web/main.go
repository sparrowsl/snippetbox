package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	address := flag.String("addr", ":5000", "HTTP Network Address")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", viewSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Printf("Starting server on port %s", *address)
	err := http.ListenAndServe(*address, mux)
	log.Fatal(err)
}
