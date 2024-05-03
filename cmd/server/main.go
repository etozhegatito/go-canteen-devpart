package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go-canteen-devpart/auth"
	"go-canteen-devpart/db"
	"net/http"
)

func main() {
	//Подключаемся к postgreSQL и иницализируем модели
	db.ConnectDatabase()

	// Инициализация хранилища для сессий
	store := cookie.NewStore([]byte("secret"))
	gateway := gin.Default()
	gateway.LoadHTMLGlob("web/templates/*")
	gateway.Use(sessions.Sessions("mysession", store))
	gateway.Use(cors.Default())

	// Обеспечение доступа к статическим файлам
	gateway.Static("/web", "./web")

	// Обработчики маршрутов для функций приложения
	// Просмотр и редактирование/обновление всех блюд
	gateway.GET("/dishes", db.GetDishes)
	gateway.POST("/dishes", db.CreateDish)
	gateway.PUT("/dishes/:id", db.UpdateDish)
	gateway.DELETE("/dishes/:id", db.DeleteDish)

	// Сделать заказ
	gateway.POST("/orders", db.CreateOrder)
	gateway.GET("/orders", func(requests *gin.Context) {
		requests.File("web/templates/order.html")
	})

	// Регистрация и вход
	gateway.POST("/signup", auth.SignUp)
	gateway.GET("/signup", func(requests *gin.Context) {
		requests.File("web/templates/signup.html")
	})
	gateway.POST("/signin", auth.SignIn)
	gateway.GET("/signin", func(requests *gin.Context) {
		requests.File("web/templates/signin.html")
	})

	// Выход
	gateway.GET("/logout", func(requests *gin.Context) {
		session := sessions.Default(requests)
		session.Delete("user_id")
		session.Save()
		requests.Redirect(http.StatusFound, "/signin")
	})

	// Проверка административных прав и панель управление блюд
	gateway.GET("/adminPage", AdminPageAuth)
	gateway.GET("/analytics", db.GetAnalytics)

	// Дашборд
	gateway.GET("/dashboard", DashboardAuth)

	// Запуск сервера
	gateway.Run(":8080")
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
