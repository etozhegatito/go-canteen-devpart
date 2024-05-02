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
	db.ConnectDatabase()

	// Инициализация хранилища для сессий
	store := cookie.NewStore([]byte("secret"))

	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")
	router.Use(sessions.Sessions("mysession", store))
	router.Use(cors.Default())

	// Обеспечение доступа к статическим файлам
	router.Static("/web", "./web")

	// Обработчики маршрутов для основных функций приложения
	router.GET("/dishes", db.GetDishes)
	router.POST("/dishes", db.CreateDish)
	router.PUT("/dishes/:id", db.UpdateDish)
	router.DELETE("/dishes/:id", db.DeleteDish)

	router.POST("/orders", db.CreateOrder)
	router.GET("/orders", func(c *gin.Context) {
		c.File("web/templates/order.html")
	})

	// Регистрация и вход
	router.POST("/signup", auth.SignUp)
	router.GET("/signup", func(c *gin.Context) {
		c.File("web/templates/signup.html")
	})
	router.POST("/signin", auth.SignIn)
	router.GET("/signin", func(c *gin.Context) {
		c.File("web/templates/signin.html")
	})

	// Выход
	router.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("user_id")
		session.Save()
		c.Redirect(http.StatusFound, "/signin")
	})

	// Панель администратора с проверкой административных прав
	router.GET("/adminPage", AdminPageHandler)
	router.GET("/dashboard", DashboardHandler)
	router.GET("/analytics", db.GetAnalytics)
	// Запуск сервера
	router.Run(":8080")
}

func AdminPageHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	user, err := db.GetUserByID(userID.(uint))
	if err != nil || !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.HTML(http.StatusOK, "adminPage.html", gin.H{"UserName": user.Name})
}

func DashboardHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	user, err := db.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{"UserName": user.Name})
}
