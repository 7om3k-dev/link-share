package main

import (
	"log"
	"net/http"

	"github.com/7om3k/link-share/auth-server/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /authorise", handlers.Authorise)
	mux.HandleFunc("POST /approve", handlers.Approve)

	srv := http.Server{
		Addr:    "localhost:3000",
		Handler: mux,
	}

	log.Print("starting auth server")
	log.Fatal(srv.ListenAndServe())
}
