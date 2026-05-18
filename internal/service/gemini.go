package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/genai"
)

const (
	GeminiModel = "gemini-2.5-flash-lite"
	AITimeout   = 30 * time.Second
)

var (
	ErrAITimeout     = errors.New("ai request timeout")
	ErrAIUnavailable = errors.New("ai service unavailable")
)

// GeminiClient is a shared wrapper around the Google Generative AI SDK.
// Create once at startup and pass to all handlers.
type GeminiClient struct {
	client *genai.Client
}

// NewGeminiClient creates a GeminiClient using the provided API key.
func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}
	return &GeminiClient{client: client}, nil
}

// Generate sends userText to Gemini with the given systemPrompt.
// The request is cancelled after AITimeout (30 s).
// Returns ErrAITimeout on deadline exceeded, ErrAIUnavailable on any other error.
func (g *GeminiClient) Generate(systemPrompt, userText string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), AITimeout)
	defer cancel()

	result, err := g.client.Models.GenerateContent(
		ctx,
		GeminiModel,
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
