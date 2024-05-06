package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-canteen-devpart/db"
	"go-canteen-devpart/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// SignUp для нового юзера
func SignUp(requests *gin.Context) {
	var user models.User

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
	var creds models.Credentials

	//проверяем данные пользователся на валидность
	if err := requests.ShouldBindJSON(&creds); err != nil {
		requests.JSON(http.StatusBadRequest, gin.H{"error": "Все по пизде пошло, давай нормально"})
		return
	}

	//Чекаем есть вообще такой юзер
	log.Println("Credentials received:", creds)
	db.CheckUser(creds, requests)

}

func AdminPageAuth(requests *gin.Context) {
	session := sessions.Default(requests)
	userID := session.Get("user_id")

	//проверка наличие авторизаций
	if userID == nil {
		requests.Redirect(http.StatusFound, "/signin")
		return
	}

	//проверка через базу данных наличие админки
	user, err := db.GetUserByID(userID.(uint))
	if err != nil || !user.IsAdmin {
		requests.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	//доступ разрешен
	requests.HTML(http.StatusOK, "adminPage.html", gin.H{"UserName": user.Name})
}

func DashboardAuth(requests *gin.Context) {
	session := sessions.Default(requests)
	userID := session.Get("user_id")

	//проверка наличие авторизаций
	if userID == nil {
		requests.Redirect(http.StatusFound, "/signin")
		return
	}

	//Получаем данные юзера для отображение в дашборде, и обработчик ошибка на всякий случай
	user, err := db.GetUserByID(userID.(uint))
	if err != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in database"})
		return
	}

	requests.HTML(http.StatusOK, "dashboard.html", gin.H{"UserName": user.Name})
}
