package main

import (
	"log"
	"net/http"

	"github.com/umutondersu/mathapi/internal/middleware"
	"github.com/umutondersu/mathapi/internal/routes"
)

func main() {
	router := routes.NewRouter()
	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.ChainMiddleware(router, middleware.LimitMiddleware, middleware.LoggingMiddleware),
	}

	log.Println("\nStarting server on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
