package main

import (
	"github.com/gin-gonic/gin"
	"go-canteen-devpart/auth"
	"go-canteen-devpart/db"
	"log"
	"net/http"
)

type OrderRequest struct {
	UserID    uint       `json:"user_id"`
	CartItems []CartItem `json:"cart_items"`
}

type CartItem struct {
	DishID   uint `json:"dish_id"`
	Quantity int  `json:"quantity"`
}

func main() {
	db.ConnectDatabase()
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)

	log.Println("Starting server on :7070...")
	if err := http.ListenAndServe(":7070", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	router := gin.Default()
	router.POST("/orders", db.CreateOrder)
	router.Run(":7071")

	//
	//r := gin.Default() // Инициализируем Gin
	//
	//r.POST("/signup", api.SignUp) // Регистрируем роут для регистрации
	//r.POST("/signin", api.SignIn) // Регистрируем роут для входа
	//
	//// Защищенные роуты
	//authGroup := r.Group("/")
	//authGroup.Use(api.AuthMiddleware()) // Применяем middleware для аутентификации
	//{
	//	// Тут будут защищенные роуты
	//}
	//
	//log.Println("Starting server on :7070...")
	//if err := r.Run(":7070"); err != nil { // Запускаем сервер через Gin
	//	log.Fatal("Error starting server: ", err)
	//}

}
