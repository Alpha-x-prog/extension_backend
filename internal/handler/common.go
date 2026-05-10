package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Alpha-x-prog/extension_backend/internal/service"
)

// ErrorResponse is returned on all non-2xx responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeAIError maps service sentinel errors to the correct HTTP status codes.
func writeAIError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrAITimeout:
		writeJSON(w, http.StatusGatewayTimeout, ErrorResponse{Error: "request timeout"})
	case service.ErrAIUnavailable:
		writeJSON(w, http.StatusBadGateway, ErrorResponse{Error: "AI service unavailable"})
	default:
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
