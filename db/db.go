package db

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var database *gorm.DB
var err error

func ConnectDatabase() {
	dbs := "host=localhost user=postgres dbname=postgres password=mysecretpassword port=5432"
	database, err = gorm.Open(postgres.Open(dbs), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	} else {
		log.Println("Successfully connected to database")
	}
	err := database.AutoMigrate(&User{}) // AutoMigrate для модели User
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
	database.AutoMigrate(&User{}, &Dish{}, &Order{}, &OrderItem{})

	user := User{Name: "Aqa", Surname: "Aqa", Password: "123", Phone: "111", Email: "pro", Age: 18, IsAdmin: true}
	menu := Dish{Name: "salad", Price: 50, Description: "prosto", Weight: 11, InStock: true}

	database.Create(&menu)
	database.Create(&user)
	var users []User
	result := database.Where("name = ?", "Aka").First(&users)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("No user found with the name Alice.")
		} else {
			log.Fatal("Error searching for user:", result.Error)
		}
	} else {
		fmt.Printf("User found: %+v\n", user)
	}

	var menuse []Dish
	resultt := database.Where("name = ?", "salad").First(&menuse)
	if resultt.Error != nil {
		if resultt.Error == gorm.ErrRecordNotFound {
			fmt.Println("No user found with the name Alice.")
		} else {
			log.Fatal("Error searching for user:", resultt.Error)
		}
	} else {
		fmt.Printf("User found: %+v\n", menuse)
	}

}

func CreateOrder(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var totalSum float64
	for _, item := range req.CartItems {
		var dish Dish
		if result := database.First(&dish, item.DishID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dish not found"})
			return
		}
		totalSum += float64(item.Quantity) * dish.Price
	}

	newOrder := Order{UserID: req.UserID, TotalSum: totalSum}
	if result := database.Create(&newOrder); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	for _, item := range req.CartItems {
		newOrderItem := OrderItem{OrderID: newOrder.ID, DishID: item.DishID, Quantity: item.Quantity}
		database.Create(&newOrderItem)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order_id": newOrder.ID})
}

func GetDishes(c *gin.Context) {
	var dishes []Dish
	if result := database.Find(&dishes); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, dishes)
}

func CreateUser(user User, c *gin.Context) {
	if err := database.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}
}

func CheckUser(creds Credentials, c *gin.Context) {
	var user User
	log.Println("Checking user with email:", creds.Email)

	if result := database.Where("email = ?", creds.Email).First(&user); result.Error != nil {
		log.Println("User not found:", result.Error)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Println("Password comparison failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := createToken(user.ID)
	if err != nil {
		log.Println("Token creation failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	log.Println("Token created successfully:", token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func createToken(userID uint) (string, error) {
	var jwtKey = []byte("your_secret_key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
