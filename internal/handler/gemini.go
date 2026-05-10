package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/genai"
)

const (
	geminiModel = "gemini-2.5-flash"
	aiTimeout   = 30 * time.Second
)

var (
	ErrAITimeout     = errors.New("ai request timeout")
	ErrAIUnavailable = errors.New("ai service unavailable")
)

// GeminiClient wraps the Google Generative AI client and is shared across handlers.
type GeminiClient struct {
	client *genai.Client
}

// NewGeminiClient creates a single shared Gemini client. Call once at startup.
func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}
	return &GeminiClient{client: client}, nil
}

// Generate sends userText to Gemini with the given systemPrompt.
// Returns ErrAITimeout on 30-second timeout, ErrAIUnavailable on any other error.
func (g *GeminiClient) Generate(systemPrompt, userText string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), aiTimeout)
	defer cancel()

	result, err := g.client.Models.GenerateContent(
		ctx,
		geminiModel,
		genai.Text(userText),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: systemPrompt}},
			},
		},
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Gemini timeout: %v", err)
			return "", ErrAITimeout
		}
		log.Printf("Gemini error: %v", err)
		return "", ErrAIUnavailable
	}

	return result.Text(), nil
}
