package auth

import (
	"github.com/gin-gonic/gin"
	"go-canteen-devpart/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var jwtKey = []byte("your_secret_key") // Replace with a secret key in a secure way

// SignUp handles registration of a new user
func SignUp(c *gin.Context) {
	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data!", "details": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Создание пользователя в базе данных
	db.CreateUser(user, c)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}

// SignIn handles user login
func SignIn(c *gin.Context) {
	var creds db.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data!"})
		return
	}

	log.Println("Credentials received:", creds)
	db.CheckUser(creds, c)
}
