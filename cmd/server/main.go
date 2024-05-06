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
	gateway.Static("/web/templates/", "./web")

	// Обработчики маршрутов для функций приложения
	// Просмотр и редактирование/обновление всех блюд
	gateway.GET("/dishes", db.GetDishes)
	gateway.POST("/dishes", db.CreateDish)
	gateway.PUT("/dishes/:id", db.UpdateDish)
	gateway.DELETE("/dishes/:id", db.DeleteDish)

	// Сделать заказ
	gateway.POST("/orders", db.CreateOrder)
	gateway.GET("/orders", func(c *gin.Context) {
		c.HTML(http.StatusOK, "order.html", gin.H{})
	})

	// Регистрация и вход
	gateway.POST("/signup", auth.SignUp)
	gateway.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	})

	gateway.POST("/signin", auth.SignIn)
	gateway.GET("/signin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signin.html", gin.H{})
	})
	// Выход
	gateway.GET("/logout", func(requests *gin.Context) {
		session := sessions.Default(requests)
		session.Delete("user_id")
		session.Save()
		requests.Redirect(http.StatusFound, "/signin")
	})

	// Проверка административных прав и панель управление блюд
	gateway.GET("/adminPage", auth.AdminPageAuth)
	gateway.GET("/analytics", db.OrdersHandler)

	// Дашборд
	gateway.GET("/", auth.DashboardAuth)

	// Запуск сервера
	gateway.Run(":8080")
}
