package api

import (
	"net/http"
	"strings"
)

// AuthMiddleware проверяет наличие и валидность токена аутентификации
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Authorization header must be in the format 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		// Допустим, здесь должна быть проверка токена
		token := headerParts[1]
		if token != "expected_token" { // Замените на реальную проверку токена
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Токен валиден, продолжаем обработку запроса
		next.ServeHTTP(w, r)
	})
}
