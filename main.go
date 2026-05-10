package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Alpha-x-prog/extension_backend/internal/handler"
)

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	gc, err := handler.NewGeminiClient(apiKey)
	if err != nil {
		log.Fatalf("failed to create Gemini client: %v", err)
	}

	http.HandleFunc("/explain", handler.NewExplainHandler(gc))
	http.HandleFunc("/summary", handler.NewSummaryHandler(gc))
	http.HandleFunc("/glossary", handler.NewGlossaryHandler(gc))

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
