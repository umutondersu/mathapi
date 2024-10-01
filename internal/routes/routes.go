package routes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /echo/{id}", echoHandler)
	router.HandleFunc("POST /add", addHandler)

	return router
}

type Operation struct {
	Number1 int `json:"number1"`
	Number2 int `json:"number2"`
}

type OperationResult struct {
	Result int `json:"result"`
}

type SumOperation struct {
	Numbers []int `json:"numbers"`
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	id := r.PathValue("id")
	w.Write([]byte("Recieved request for item: " + id))
	logger.Info(fmt.Sprintf("Received request for item: %s", id))
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var numbers Operation

	err := json.NewDecoder(r.Body).Decode(&numbers)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		logger.Error("Failed to Decode response", slog.String("error", err.Error()))
		return
	}

	result := OperationResult{Result: numbers.Number1 + numbers.Number2}
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("Failed to marshal response", slog.String("error", err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	w.Header().Set("Content-Type", "application/json")
	logger.Info("Successfully added two numbers", slog.Int("result", result.Result))
}
