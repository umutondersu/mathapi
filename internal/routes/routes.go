package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /echo/{id}", handleEcho)
	router.HandleFunc("POST /add", handleAdd)
	router.HandleFunc("POST /substract", handleSubstract)

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

func handleEcho(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	id := r.PathValue("id")
	w.Write([]byte("Recieved request for item: " + id))
	logger.Info(fmt.Sprintf("Received request for item: %s", id))
}

func checkInput(op *Operation, w http.ResponseWriter, r *http.Request, logger *slog.Logger) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&op)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeError) {
			http.Error(w, "number1 and number2 must be numbers", http.StatusBadRequest)
			logger.Error("number1 and number2 must be numbers", slog.String("error", err.Error()))
			return errors.New("number1 and number2 must be numbers")
		} else {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			logger.Error("Failed to decode request body", slog.String("error", err.Error()))
			return errors.New("Failed to decode request body")
		}
	}
	return nil
}

func marshalResult(result OperationResult, w http.ResponseWriter, logger *slog.Logger) ([]byte, error) {
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("Failed to marshal response", slog.String("error", err.Error()))
		return nil, errors.New("Failed to marshal response")
	}
	return response, nil
}

func prepareResponse(w http.ResponseWriter, response []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	w.Header().Set("Content-Type", "application/json")
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var numbers Operation

	if err := checkInput(&numbers, w, r, logger); err != nil {
		return
	}

	result := OperationResult{Result: numbers.Number1 + numbers.Number2}
	response, err := marshalResult(result, w, logger)
	if err != nil {
		return
	}

	prepareResponse(w, response)

	logger.Info("Successfully added two numbers", slog.Int("result", result.Result))
}

func handleSubstract(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var numbers Operation

	if err := checkInput(&numbers, w, r, logger); err != nil {
		return
	}

	result := OperationResult{Result: numbers.Number1 - numbers.Number2}
	response, err := marshalResult(result, w, logger)
	if err != nil {
		return
	}

	prepareResponse(w, response)

	logger.Info("Successfully substracted two numbers", slog.Int("result", result.Result))
}
