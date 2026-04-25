package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type ExplainRequest struct {
	Text string `json:"text"`
}

type ExplainResponse struct {
	Text string `json:"text"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func explainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	if r.Body == nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	var req ExplainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	if len(req.Text) == 0 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	if len(req.Text) > 1000 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "text exceeds maximum length of 1000 characters"})
		return
	}

	writeJSON(w, http.StatusOK, ExplainResponse{Text: req.Text})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/explain", explainHandler)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
