//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/ollama/ollama/api"
)

func TestOllamaConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	req := api.ChatRequest{
		Model: smol,
		Messages: []api.Message{
			{
				Role:    "user",
				Content: "Say hello",
			},
		},
		Stream: &stream,
		Options: map[string]any{
			"temperature": 0,
			"seed":        123,
		},
	}
	ChatTestHelper(ctx, t, req, "hello")
}
