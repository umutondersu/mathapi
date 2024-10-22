package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestOperation struct {
	Number1 any `json:"number1"`
	Number2 any `json:"number2"`
}

type TestSumOperation struct {
	Numbers []any `json:"numbers"`
}

type TestSumOperationInvalidkey struct {
	Numbs []any `json:"numbs"`
}

type TestOperationinvalidkeys struct {
	Numb1 any `json:"numb1"`
	Numb2 any `json:"numb2"`
}

var (
	k InvalidKeysError
	v BOValuesError
)

func TestEchoHandler(t *testing.T) {
	tests := []struct {
		request          int
		expectedResponse string
	}{
		{1, "Received request for item: 1"},
		{2, "Received request for item: 2"},
		{-1, "Received request for item: -1"},
		{0, "Received request for item: 0"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", fmt.Sprintf("/echo/%d", tt.request), nil)
		w := httptest.NewRecorder()

		handleEcho(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(body) != tt.expectedResponse {
			t.Fatalf("Expected response to be %s, got %s", tt.expectedResponse, string(body))
		}
	}
}

func TestAddHandler(t *testing.T) {
	tests := []struct {
		data     any
		expected any
	}{
		{TestOperation{1, 2}, 3},
		{TestOperation{2, 3}, 5},
		{TestOperation{0, 0}, 0},
		{TestOperation{-1, 1}, 0},
		{TestOperation{-1, -1}, -2},
		{TestOperation{5, string('b')}, v.Error()},
		{TestOperation{string('a'), string('b')}, v.Error()},
		{TestOperation{string('a'), 5}, v.Error()},
		{TestOperationinvalidkeys{1, 2}, k.Error()},
		{TestOperationinvalidkeys{string('a'), 2}, k.Error()},
		{TestOperationinvalidkeys{string('b'), string('a')}, k.Error()},
		{"This is a string request", k.Error()},
	}
	for _, tt := range tests {
		// Encode the data to JSON
		jsonData, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Create a new HTTP request
		req := httptest.NewRequest("POST", "/add", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		// Call the handler with the request
		handleAdd(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		// Check for non-JSON response for invalid types of value (Bad Request)
		if resp.StatusCode == http.StatusBadRequest {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			if strings.TrimSpace(string(body)) != tt.expected {
				t.Fatalf("For Bad Request Expected response to be %s, got %s", tt.expected, string(body))
			}
			return
		}

		// Read and Decode the response JSON
		var result OperationResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response to OperationResult: %v", err)
		}
		if result.Result != tt.expected {
			t.Fatalf("Expected response to be %d, got %d", tt.expected, result.Result)
		}
	}
}

func TestSubtractdHandler(t *testing.T) {
	tests := []struct {
		data     any
		expected any
	}{
		{TestOperation{3, 1}, 2},
		{TestOperation{0, 0}, 0},
		{TestOperation{3, 5}, -2},
		{TestOperation{-1, 1}, -2},
		{TestOperation{-1, -1}, 0},
		{TestOperation{5, string('b')}, v.Error()},
		{TestOperation{string('a'), string('b')}, v.Error()},
		{TestOperation{string('a'), 5}, v.Error()},
		{TestOperationinvalidkeys{1, 2}, k.Error()},
		{TestOperationinvalidkeys{string('a'), 2}, k.Error()},
		{TestOperationinvalidkeys{string('b'), string('a')}, k.Error()},
		{"This is a string request", k.Error()},
	}
	for _, tt := range tests {
		// Encode the data to JSON
		jsonData, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Create a new HTTP request
		req := httptest.NewRequest("POST", "/subtract", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		// Call the handler with the request
		handleSubtract(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		// Check for non-JSON response for invalid types of value (Bad Request)
		if resp.StatusCode == http.StatusBadRequest {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			if strings.TrimSpace(string(body)) != tt.expected {
				t.Fatalf("For Bad Request Expected response to be %s, got %s", tt.expected, string(body))
			}
			return
		}

		// Read and Decode the response JSON
		var result OperationResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response to OperationResult: %v", err)
		}
		if result.Result != tt.expected {
			t.Fatalf("Expected response to be %d, got %d", tt.expected, result.Result)
		}
	}
}

func TestMultiplyHandler(t *testing.T) {
	tests := []struct {
		data     any
		expected any
	}{
		{TestOperation{3, 1}, 3},
		{TestOperation{0, 0}, 0},
		{TestOperation{3, 5}, 15},
		{TestOperation{-1, 1}, -1},
		{TestOperation{-1, -1}, 1},
		{TestOperation{5, string('b')}, v.Error()},
		{TestOperation{string('a'), string('b')}, v.Error()},
		{TestOperation{string('a'), 5}, v.Error()},
		{TestOperationinvalidkeys{1, 2}, k.Error()},
		{TestOperationinvalidkeys{string('a'), 2}, k.Error()},
		{TestOperationinvalidkeys{string('b'), string('a')}, k.Error()},
		{"This is a string request", k.Error()},
	}
	for _, tt := range tests {
		// Encode the data to JSON
		jsonData, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Create a new HTTP request
		req := httptest.NewRequest("POST", "/multiply", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		// Call the handler with the request
		handleMultiply(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		// Check for non-JSON response for invalid types of value (Bad Request)
		if resp.StatusCode == http.StatusBadRequest {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			if strings.TrimSpace(string(body)) != tt.expected {
				t.Fatalf("For Bad Request Expected response to be %s, got %s", tt.expected, string(body))
			}
			return
		}

		// Read and Decode the response JSON
		var result OperationResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response to OperationResult: %v", err)
		}
		if result.Result != tt.expected {
			t.Fatalf("Expected response to be %d, got %d", tt.expected, result.Result)
		}
	}
}

func TestDivideHandler(t *testing.T) {
	var divisionbyzero DivisionbyZeroError
	tests := []struct {
		data     any
		expected any
	}{
		{TestOperation{3, 1}, 3},
		{TestOperation{0, 0}, divisionbyzero.Error()},
		{TestOperation{5, 0}, divisionbyzero.Error()},
		{TestOperation{3, 5}, 15},
		{TestOperation{-1, 1}, -1},
		{TestOperation{-1, -1}, 1},
		{TestOperation{5, string('b')}, v.Error()},
		{TestOperation{string('a'), string('b')}, v.Error()},
		{TestOperation{string('a'), 5}, v.Error()},
		{TestOperationinvalidkeys{1, 2}, k.Error()},
		{TestOperationinvalidkeys{string('a'), 2}, k.Error()},
		{TestOperationinvalidkeys{string('b'), string('a')}, k.Error()},
		{"This is a string request", k.Error()},
	}
	for _, tt := range tests {
		// Encode the data to JSON
		jsonData, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Create a new HTTP request
		req := httptest.NewRequest("POST", "/divide", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		// Call the handler with the request
		handleDivide(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		// Check for non-JSON response for invalid types of value (Bad Request)
		if resp.StatusCode == http.StatusBadRequest {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			if strings.TrimSpace(string(body)) != tt.expected {
				t.Fatalf("For Bad Request Expected response to be %s, got %s", tt.expected, string(body))
			}
			return
		}

		// Read and Decode the response JSON
		var result OperationResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response to OperationResult: %v", err)
		}
		if result.Result != tt.expected {
			t.Fatalf("Expected response to be %d, got %d", tt.expected, result.Result)
		}
	}
}

func TestSumHandler(t *testing.T) {
	var v SOValueError
	tests := []struct {
		data     any
		expected any
	}{
		{TestSumOperation{[]any{1, 2, 3, 4, 5}}, 15},
		{TestSumOperation{[]any{1, -1, 0, 1, -1}}, 0},
		{TestSumOperation{[]any{1, 2, 3, -4, -5}}, -3},
		{TestSumOperation{[]any{5, string('b')}}, v.Error()},
		{TestSumOperationInvalidkey{[]any{string('a'), 5}}, k.Error()},
		{TestSumOperationInvalidkey{[]any{1, 2}}, k.Error()},
		{"This is a string request", k.Error()},
	}
	for i, tt := range tests {
		// Encode the data to JSON
		jsonData, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Create a new HTTP request
		req := httptest.NewRequest("POST", "/sum", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()

		// Call the handler with the request
		handleSum(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		// Check for non-JSON response for invalid types of value (Bad Request)
		if resp.StatusCode == http.StatusBadRequest {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			if strings.TrimSpace(string(body)) != tt.expected {
				t.Fatalf("At test case %v, For Bad Request Expected response to be %s, got %s", i+1, tt.expected, string(body))
			}
			return
		}

		// Read and Decode the response JSON
		var result OperationResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response to OperationResult: %v", err)
		}
		if result.Result != tt.expected {
			t.Fatalf("Expected response to be %d, got %d", tt.expected, result.Result)
		}
	}
}
