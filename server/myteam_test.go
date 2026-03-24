package server

import (
    "testing"
    "net/http"
)

func TestOllamaServerIsReachable(t *testing.T) {
    resp, err := http.Get("http://localhost:11434")
    if err != nil {
        t.Skip("Ollama not running, skipping test")
    }
    if resp.StatusCode != 200 {
        t.Errorf("Expected status 200, got %d", resp.StatusCode)
    }
}