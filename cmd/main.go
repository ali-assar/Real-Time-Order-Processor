package main

import (
	"log"
	"net/http"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/handler"
)

func main() {

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("API listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
