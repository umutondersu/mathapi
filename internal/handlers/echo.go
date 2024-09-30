package handlers

import (
	"log"
	"net/http"
)

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request %s %s", r.Method, r.URL.Path)
	id := r.PathValue("id")
	w.Write([]byte("Recieved request for item: " + id))
}
