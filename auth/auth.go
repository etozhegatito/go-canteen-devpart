package auth

import (
	"encoding/json"
	"net/http"
)

// Credentials структура для данных аутентификации пользователя.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterHandler обрабатывает запросы на регистрацию новых пользователей.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Здесь добавьте логику добавления нового пользователя в базу данных
	// Это мог бы быть вызов функции, который вы реализуете в internal/db
}

// LoginHandler обрабатывает запросы на вход от пользователей.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Здесь добавьте логику проверки данных пользователя
	// Это мог бы быть вызов функции, который вы реализуете в internal/db
}
