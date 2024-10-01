package routes

import (
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
