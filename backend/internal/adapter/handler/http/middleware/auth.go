package middleware

import (
	"net/http"
	"strings"
)

// Mock credentials
const (
	mockAPIKey   = "mock-api-key"
	mockJWTToken = "mock-jwt-token"
)

// Authenticator is a simple middleware to check for a mock API key or JWT token.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for API Key
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == mockAPIKey {
			next.ServeHTTP(w, r)
			return
		}

		// Check for Bearer Token
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == mockJWTToken {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

