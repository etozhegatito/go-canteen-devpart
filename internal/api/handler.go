package api

import (
	"encoding/json"
	"net/http"
)

// RegisterHandler обрабатывает регистрацию новых пользователей
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Здесь должен быть код для регистрации пользователя
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User registered successfully")
}

// LoginHandler обрабатывает аутентификацию пользователей
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string
		Password string
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if result := db.DB.Where("email = ?", creds.Email).First(&user); result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !VerifyPassword(user.Password, creds.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User logged in successfully")
}
