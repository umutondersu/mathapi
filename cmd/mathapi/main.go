package main

import (
	"log"
	"net/http"

	"github.com/umutondersu/mathapi/internal/handlers"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/echo/{id}", handlers.EchoHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("\nStarting server on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
