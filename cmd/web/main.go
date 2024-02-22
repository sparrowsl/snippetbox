package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	address := flag.String("addr", ":5000", "HTTP Network Address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", viewSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	server := &http.Server{
		Addr:     *address,
		Handler:  mux,
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on port %s", *address)
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
