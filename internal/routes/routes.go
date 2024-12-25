package routes

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/umutondersu/mathapi/internal/middleware"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /echo/{id}", handleEcho)
	router.HandleFunc("POST /add", handleAdd)
	router.HandleFunc("POST /subtract", handleSubtract)
	router.HandleFunc("POST /multiply", handleMultiply)
	router.HandleFunc("POST /divide", handleDivide)
	router.HandleFunc("POST /sum", handleSum)

	return router
}

type BasicOperation struct {
	Number1 int `json:"number1"`
	Number2 int `json:"number2"`
}

type OperationResult struct {
	Result int `json:"result"`
}

type SumOperation struct {
	Numbers []int `json:"numbers"`
}

type Operation int

const (
	ADD = Operation(iota)
	SUBTRACT
	MULTIPLY
	DIVIDE
	SUM
)

func handleEcho(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	// id := r.PathValue("id")  // Pathvalue does not work with httptest
	id := r.URL.Path[len("/echo/"):]

	_, err := w.Write([]byte("Received request for item: " + id))
	if err != nil {
		logger.Error("Failed to write response", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	logger.Info("Received request for item:", slog.String("id", id))
}

func decodeInput(w http.ResponseWriter, r *http.Request, logger *slog.Logger, op interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(op)
	if err != nil {
		errorMsg := InvalidKeysError{}.Error()
		var unmarshalTypeError *json.UnmarshalTypeError

		if errors.As(err, &unmarshalTypeError) {
			switch op.(type) {
			case *SumOperation:
				errorMsg = SOValueError{}.Error()
			case *BasicOperation:
				errorMsg = BOValuesError{}.Error()
			}
		}

		http.Error(w, errorMsg, http.StatusBadRequest)
		logger.Error(errorMsg, slog.String("error", err.Error()))
		return errors.New(errorMsg)
	}
	return nil
}

func prepareResponse(result OperationResult, w http.ResponseWriter, logger *slog.Logger) error {
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

func operationHandler(w http.ResponseWriter, r *http.Request, logger *slog.Logger, op Operation) (OperationResult, error) {
	var numbers BasicOperation
	var sumNumbers SumOperation
	var result OperationResult

	if op == SUM {
		if err := decodeInput(w, r, logger, &sumNumbers); err != nil {
			return result, err
		}
	} else {
		if err := decodeInput(w, r, logger, &numbers); err != nil {
			return result, err
		}
	}

	switch op {
	case ADD:
		result = OperationResult{Result: numbers.Number1 + numbers.Number2}
	case SUBTRACT:
		result = OperationResult{Result: numbers.Number1 - numbers.Number2}
	case MULTIPLY:
		result = OperationResult{Result: numbers.Number1 * numbers.Number2}
	case DIVIDE:
		if numbers.Number2 == 0 {
			var divisionbyzero DivisionbyZeroError
			logger.Error(divisionbyzero.Error(), slog.Int("number1", numbers.Number1), slog.Int("number2", numbers.Number2))
			http.Error(w, divisionbyzero.Error(), http.StatusBadRequest)
			return result, errors.New(divisionbyzero.Error())
		}
		result = OperationResult{Result: numbers.Number1 / numbers.Number2}
	case SUM:
		sum := 0
		for _, num := range sumNumbers.Numbers {
			sum += num
		}
		result = OperationResult{Result: sum}
	default:
		return result, errors.New("Invalid Operation")
	}

	if err := prepareResponse(result, w, logger); err != nil {
		return result, err
	}

	return result, nil
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	result, err := operationHandler(w, r, logger, ADD)
	if err != nil {
		return
	}
	logger.Info("Successfully added two numbers", slog.Int("result", result.Result))
}

func handleSubtract(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	result, err := operationHandler(w, r, logger, SUBTRACT)
	if err != nil {
		return
	}
	logger.Info("Successfully subtracted two numbers", slog.Int("result", result.Result))
}

func handleMultiply(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	result, err := operationHandler(w, r, logger, MULTIPLY)
	if err != nil {
		return
	}
	logger.Info("Successfully multiplied two numbers", slog.Int("result", result.Result))
}

func handleDivide(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	result, err := operationHandler(w, r, logger, DIVIDE)
	if err != nil {
		return
	}
	logger.Info("Successfully divided two numbers", slog.Int("result", result.Result))
}

func handleSum(w http.ResponseWriter, r *http.Request) {
	logger := middleware.GetLogger(r)
	result, err := operationHandler(w, r, logger, SUM)
	if err != nil {
		return
	}
	logger.Info("Successfully summed up Numbers", slog.Int("result", result.Result))
}
