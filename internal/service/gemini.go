package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/genai"
)

const (
	defaultGeminiModel = "gemini-2.5-flash-lite"
	AITimeout          = 30 * time.Second
)

var (
	ErrAITimeout     = errors.New("ai request timeout")
	ErrAIUnavailable = errors.New("ai service unavailable")
)

// GeminiClient is a shared wrapper around the Google Generative AI SDK.
// Create once at startup and pass to all handlers.
type GeminiClient struct {
	client *genai.Client
	model  string
}

// NewGeminiClient creates a GeminiClient using the provided API key.
// Model is read from GEMINI_MODEL env var, falls back to gemini-2.5-flash-lite.
func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = defaultGeminiModel
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}

	log.Printf("Using Gemini model: %s", model)
	return &GeminiClient{client: client, model: model}, nil
}

// Generate sends userText to Gemini with the given systemPrompt.
// The request is cancelled after AITimeout (30 s).
// Returns ErrAITimeout on deadline exceeded, ErrAIUnavailable on any other error.
func (g *GeminiClient) Generate(systemPrompt, userText string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), AITimeout)
	defer cancel()

	result, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		genai.Text(userText),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: systemPrompt}},
			},
		},
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("gemini timeout: %v", err)
			return "", ErrAITimeout
		}
		log.Printf("gemini error: %v", err)
		return "", ErrAIUnavailable
	}

	return result.Text(), nil
}
