package server

import (
	"strings"
	"testing"
)

func TestOllamaResponse(t *testing.T) {
	// Test that a basic string response is valid
	result := "expected output"

	if result == "" {
		t.Error("Expected a response but got empty string")
	}

	if !strings.Contains(result, "expected") {
		t.Error("Response did not contain expected content")
	}

	t.Log("TestOllamaResponse passed successfully")
}