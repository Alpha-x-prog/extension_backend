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

// CORSMiddleware добавляет заголовки CORS ко всем ответам.
// Также обрабатывает preflight-запросы (OPTIONS), которые браузер
// отправляет перед основным запросом для проверки разрешений.
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
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
