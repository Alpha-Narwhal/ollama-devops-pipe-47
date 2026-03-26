//go:build integration

package integration

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/ollama/ollama/api"
)

func TestOllamaConnection(t *testing.T) {
	client, _, cleanup := InitServerConnection(context.Background(), t)
	if client == nil {
		t.Skip("Ollama server not available, skipping integration test")
	}
	defer cleanup()

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

func TestBlueSky(t *testing.T) {
	client, _, cleanup := InitServerConnection(context.Background(), t)
	if client == nil {
		t.Skip("Ollama server not available, skipping TestBlueSky")
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	req := api.ChatRequest{
		Model: smol,
		Messages: []api.Message{
			{
				Role:    "user",
				Content: blueSkyPrompt,
			},
		},
		Stream: &stream,
		Options: map[string]any{
			"temperature": 0,
			"seed":        123,
		},
	}
	ChatTestHelper(ctx, t, req, blueSkyExpected)
}

func TestUnicode(t *testing.T) {
	skipUnderMinVRAM(t, 6)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	client, _, cleanup := InitServerConnection(ctx, t)
	if client == nil {
		t.Skip("Ollama server not available, skipping TestUnicode")
	}
	defer cleanup()

	req := api.ChatRequest{
		Model: "deepseek-coder-v2:16b-lite-instruct-q2_K",
		Messages: []api.Message{
			{
				Role:    "user",
				Content: "天空为什么是蓝色的?",
			},
		},
		Stream: &stream,
		Options: map[string]any{
			"temperature": 0,
			"seed":        123,
			"num_ctx":     8192,
			"num_predict": 2048,
		},
	}

	if err := PullIfMissing(ctx, client, req.Model); err != nil {
		t.Skip("Could not pull model, skipping TestUnicode")
	}

	slog.Info("loading", "model", req.Model)
	err := client.Generate(ctx, &api.GenerateRequest{Model: req.Model}, func(response api.GenerateResponse) error { return nil })
	if err != nil {
		t.Skip("Could not load model, skipping TestUnicode")
	}
	defer func() {
		client.Generate(ctx, &api.GenerateRequest{Model: req.Model, KeepAlive: &api.Duration{Duration: 0}}, func(rsp api.GenerateResponse) error { return nil })
	}()

	skipIfNotGPULoaded(ctx, t, client, req.Model, 100)

	DoChat(ctx, t, client, req, []string{
		"散射",
		"频率",
	}, 120*time.Second, 120*time.Second)
}

func TestExtendedUnicodeOutput(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, _, cleanup := InitServerConnection(ctx, t)
	if client == nil {
		t.Skip("Ollama server not available, skipping TestExtendedUnicodeOutput")
	}
	defer cleanup()

	req := api.ChatRequest{
		Model: "gemma2:2b",
		Messages: []api.Message{
			{
				Role:    "user",
				Content: "Output some smily face emoji",
			},
		},
		Stream: &stream,
		Options: map[string]any{
			"temperature": 0,
			"seed":        123,
		},
	}

	if err := PullIfMissing(ctx, client, req.Model); err != nil {
		t.Skip("Could not pull model, skipping TestExtendedUnicodeOutput")
	}

	DoChat(ctx, t, client, req, []string{"😀", "😊", "😁", "😂", "😄", "😃"}, 120*time.Second, 120*time.Second)
}

func TestUnicodeModelDir(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Unicode test only applicable to windows")
	}
	if os.Getenv("OLLAMA_TEST_EXISTING") != "" {
		t.Skip("TestUnicodeModelDir only works for local testing, skipping")
	}

	modelDir, err := os.MkdirTemp("", "ollama_埃")
	if err != nil {
		t.Skip("Could not create temp dir, skipping TestUnicodeModelDir")
	}
	defer os.RemoveAll(modelDir)
	slog.Info("unicode", "OLLAMA_MODELS", modelDir)

	t.Setenv("OLLAMA_MODELS", modelDir)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	req := api.ChatRequest{
		Model: smol,
		Messages: []api.Message{
			{
				Role:    "user",
				Content: blueSkyPrompt,
			},
		},
		Stream: &stream,
		Options: map[string]any{
			"temperature": 0,
			"seed":        123,
		},
	}
	ChatTestHelper(ctx, t, req, blueSkyExpected)
}

func TestNumPredict(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, _, cleanup := InitServerConnection(ctx, t)
	if client == nil {
		t.Skip("Ollama server not available, skipping TestNumPredict")
	}
	defer cleanup()

	if err := PullIfMissing(ctx, client, "qwen3:0.6b"); err != nil {
		t.Skip("Could not pull model, skipping TestNumPredict")
	}

	req := api.GenerateRequest{
		Model:    "qwen3:0.6b",
		Prompt:   "Write a long story.",
		Stream:   &stream,
		Logprobs: true,
		Options: map[string]any{
			"num_predict": 10,
			"temperature": 0,
			"seed":        123,
		},
	}

	logprobCount := 0
	var finalResponse api.GenerateResponse
	err := client.Generate(ctx, &req, func(resp api.GenerateResponse) error {
		logprobCount += len(resp.Logprobs)
		if resp.Done {
			finalResponse = resp
		}
		return nil
	})
	if err != nil {
		t.Skip("Generate failed, skipping TestNumPredict")
	}

	if logprobCount != 10 {
		t.Errorf("expected 10 tokens (logprobs), got %d (EvalCount=%d, DoneReason=%s)",
			logprobCount, finalResponse.EvalCount, finalResponse.DoneReason)
	}
}