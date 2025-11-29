package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	handlerFS := http.FileServer(http.Dir("."))

	mux.Handle("/", handlerFS)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
