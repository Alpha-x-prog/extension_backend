package handler

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the standard error body returned on all non-2xx responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeAIError maps GeminiClient sentinel errors to correct HTTP status codes.
func writeAIError(w http.ResponseWriter, err error) {
	switch err {
	case ErrAITimeout:
		writeJSON(w, http.StatusGatewayTimeout, ErrorResponse{Error: "request timeout"})
	case ErrAIUnavailable:
		writeJSON(w, http.StatusBadGateway, ErrorResponse{Error: "AI service unavailable"})
	default:
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
