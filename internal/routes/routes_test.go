package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

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
	type TestOperation struct {
		Number1 any `json:"number1"`
		Number2 any `json:"number2"`
	}
	tests := []struct {
		data     TestOperation
		expected int
	}{
		{TestOperation{1, 2}, 3},
		{TestOperation{2, 3}, 5},
		{TestOperation{0, 0}, 0},
		{TestOperation{-1, 1}, 0},
		{TestOperation{-1, -1}, -2},
		{TestOperation{string('a'), string('b')}, 195}, // TODO: expected must be something related to "Bad Request"
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

		// Read and Decode the response
		// TODO: Handle non-json response for invalid types of value (Bad Request)

		var result OperationResult

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if result.Result != tt.expected {
			t.Fatalf("Expected response to be %d, got %d", tt.expected, result.Result)
		}
	}
}
