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
	router.HandleFunc("POST /multiply", handleMultiply)

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

type DivisionbyZeroError struct {
	firstnumber int
}

func (d DivisionbyZeroError) Error() string {
	return fmt.Sprintf("Cannot divide %v by zero", d.firstnumber)
}

type InvalidKeysError struct{}

func (i InvalidKeysError) Error() string {
	return "Bad Request"
}

type InvalidKeyValuesError struct{}

func (i InvalidKeyValuesError) Error() string {
	return "number1 and number2 must be numbers"
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	id := r.PathValue("id")
	_, err := w.Write([]byte("Recieved request for item: " + id))
	if err != nil {
		logger.Error("Failed to write response", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	logger.Info(fmt.Sprintf("Received request for item: %s", id))
}

func checkInput(op *Operation, w http.ResponseWriter, r *http.Request, logger *slog.Logger) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&op)
	if err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		var k InvalidKeysError
		var v InvalidKeyValuesError
		if errors.As(err, &unmarshalTypeError) {
			http.Error(w, v.Error(), http.StatusBadRequest)
			logger.Error(v.Error(), slog.String("error", err.Error()))
			return errors.New(v.Error())
		} else {
			http.Error(w, k.Error(), http.StatusBadRequest)
			logger.Error("Failed to decode request body", slog.String("error", err.Error()))
			return errors.New("Failed to decode request body")
		}
	}
	return nil
}

func prepareResponse(result OperationResult, w http.ResponseWriter, logger *slog.Logger) error {
	// Marshal the result into a JSON response
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.Error("Failed to marshal response", slog.String("error", err.Error()))
		return errors.New("Failed to marshal response")
	}

	// Write the response. setting the content type and status code is done automatically
	_, err = w.Write(response)
	if err != nil {
		logger.Error("Failed to write response", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return errors.New("Failed to write response")
	}

	return nil
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var numbers Operation

	if err := checkInput(&numbers, w, r, logger); err != nil {
		return
	}

	result := OperationResult{Result: numbers.Number1 + numbers.Number2}
	if err := prepareResponse(result, w, logger); err != nil {
		return
	}

	logger.Info("Successfully added two numbers", slog.Int("result", result.Result))
}

func handleSubstract(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var numbers Operation

	if err := checkInput(&numbers, w, r, logger); err != nil {
		return
	}

	result := OperationResult{Result: numbers.Number1 - numbers.Number2}
	if err := prepareResponse(result, w, logger); err != nil {
		return
	}

	logger.Info("Successfully substracted two numbers", slog.Int("result", result.Result))
}

func handleMultiply(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var numbers Operation

	if err := checkInput(&numbers, w, r, logger); err != nil {
		return
	}

	result := OperationResult{Result: numbers.Number1 * numbers.Number2}
	if err := prepareResponse(result, w, logger); err != nil {
		return
	}

	logger.Info("Successfully multiplied two numbers", slog.Int("result", result.Result))
}
