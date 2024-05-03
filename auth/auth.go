package auth

import (
	"github.com/gin-gonic/gin"
	"go-canteen-devpart/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// SignUp для нового юзера
func SignUp(requests *gin.Context) {
	var user db.User

	//Проверка на валдиность данных
	if err := requests.ShouldBindJSON(&user); err != nil {
		requests.JSON(http.StatusBadRequest, gin.H{"error": "Wrong data, try again", "details": err.Error()})
		return
	}

	// Сразу хэшируем пароль
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": "Все по пизде пошло"})
		return
	}
	user.Password = string(hashPassword)

	// Создание юзера в базе данных
	db.CreateUser(user, requests)
	requests.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}

// SignIn для входа юзера
func SignIn(requests *gin.Context) {
	//Временный Struct для создание шаблона входа
	var creds db.Credentials

	//проверяем данные пользователся на валидность
	if err := requests.ShouldBindJSON(&creds); err != nil {
		requests.JSON(http.StatusBadRequest, gin.H{"error": "Все по пизде пошло, давай нормально"})
		return
	}

	//Чекаем есть вообще такой юзер
	log.Println("Credentials received:", creds)
	db.CheckUser(creds, requests)

}
