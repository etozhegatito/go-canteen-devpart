package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	// Создаем новый запрос с нужными параметрами.
	body := []byte(`{"username":"newuser","password":"password"}`)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Используем httptest.ResponseRecorder для записи ответов.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)

	// Вызываем наш обработчик.
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код, который вернул наш обработчик.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа.
	expected := `User registered successfully`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
