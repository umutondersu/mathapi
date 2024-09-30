package routes

import (
	"log"
	"net/http"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/echo/{id}", echoHandler)

	return router
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request %s %s", r.Method, r.URL.Path)
	id := r.PathValue("id")
	w.Write([]byte("Recieved request for item: " + id))
}
