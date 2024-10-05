package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type TestOperation struct {
	Number1 any `json:"number1"`
	Number2 any `json:"number2"`
}

type TestOperationinvalidkeys struct {
	Numb1 any `json:"numb1"`
	Numb2 any `json:"numb2"`
}

func TestEchoHandler(t *testing.T) {
	tests := []struct {
		expected int
	}{
		{1},
		{2},
		{-1},
		{0},
	}

	for _, tt := range tests {
		client := &http.Client{}
		resp, err := client.Get(fmt.Sprintf("http://localhost:8080/echo/%d", tt.expected))
		if err != nil {
			t.Fatalf("Failed to get response: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(body) != fmt.Sprintf("Recieved request for item: %d", tt.expected) {
			t.Fatalf("Expected response to be Recieved request for item: %d, got %s", tt.expected, string(body))
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
		{TestOperation{5, string('b')}, "number1 and number2 must be numbers"},
		{TestOperation{string('a'), string('b')}, "number1 and number2 must be numbers"},
		{TestOperation{string('a'), 5}, "number1 and number2 must be numbers"},
		{TestOperationinvalidkeys{1, 2}, "Bad Request"},
		{TestOperationinvalidkeys{string('a'), 2}, "Bad Request"},
		{TestOperationinvalidkeys{string('b'), string('a')}, "Bad Request"},
		{"This is a string request", "Bad Request"},
	}
	for _, tt := range tests {
		client := &http.Client{}

		// Encode the data to JSON
		jsonData, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "http://localhost:8080/add", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send the request usung the HTTP client
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer req.Body.Close()

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
			t.Fatalf("Failed to decode response to TestOperationResult: %v", err)
		}
		if result.Result != tt.expected {
			t.Fatalf("Expected response to be %d, got %d", tt.expected, result.Result)
		}
	}
}

func TestSubstractdHandler(t *testing.T) {
	tests := []struct {
		data     any
		expected any
	}{
		{TestOperation{3, 1}, 2},
		{TestOperation{0, 0}, 0},
		{TestOperation{3, 5}, -2},
		{TestOperation{-1, 1}, -2},
		{TestOperation{-1, -1}, 0},
		{TestOperation{5, string('b')}, "number1 and number2 must be numbers"},
		{TestOperation{string('a'), string('b')}, "number1 and number2 must be numbers"},
		{TestOperation{string('a'), 5}, "number1 and number2 must be numbers"},
		{TestOperationinvalidkeys{1, 2}, "Bad Request"},
		{TestOperationinvalidkeys{string('a'), 2}, "Bad Request"},
		{TestOperationinvalidkeys{string('b'), string('a')}, "Bad Request"},
		{"This is a string request", "Bad Request"},
	}
	for _, tt := range tests {
		client := &http.Client{}

		// Encode the data to JSON
		jsonData, err := json.Marshal(tt.data)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		// Create a new HTTP request
		req, err := http.NewRequest("POST", "http://localhost:8080/substract", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send the request usung the HTTP client
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer req.Body.Close()

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
			t.Fatalf("Failed to decode response to TestOperationResult: %v", err)
		}
		if result.Result != tt.expected {
			t.Fatalf("Expected response to be %d, got %d", tt.expected, result.Result)
		}
	}
}
