package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-canteen-devpart/auth"
	"go-canteen-devpart/db"
	"log"
)

func main() {
	db.ConnectDatabase()

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/dishes", db.GetDishes)
	router.POST("/orders", db.CreateOrder)
	router.GET("/orders", func(c *gin.Context) {
		c.File("web/templates/order.html")
	})

	router.POST("/signup", auth.SignUp) // SignUp - функция обработчик для регистрации

	// Обеспечение доступа к HTML-страницам через пути /web/signup.html и /web/signin.html
	router.GET("/signup", func(c *gin.Context) {
		c.File("web/templates/signup.html")
	})

	router.POST("/signin", auth.SignIn) // SignUp - функция обработчик для регистрации

	// Обеспечение доступа к HTML-страницам через пути /web/signup.html и /web/signin.html
	router.GET("/signin", func(c *gin.Context) {
		c.File("web/templates/signin.html")
	})
	router.NoRoute(func(c *gin.Context) {
		log.Printf("Path: %s not found", c.Request.URL.Path)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	router.Run(":8080")
}
