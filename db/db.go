package db

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var db *gorm.DB

func ConnectDatabase() {
	dbs := "host=localhost user=postgres dbname=postgres password=mysecretpassword port=5432"
	db, err := gorm.Open(postgres.Open(dbs), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	} else {
		log.Println("Successfully connected to database")
	}

	db.AutoMigrate(&User{}, &Dish{}, Order{}, OrderItem{})
	dateStr := "1990-12-31"

	// Преобразование строки в time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	user := User{Name: "Aka", Surname: "Aka", Password: "123", Phone: "111", Email: "profaka", DateOfBirth: date, IsAdmin: true}
	menu := Dish{Name: "salad", Price: 50, Description: "prosto", Weight: 11, InStock: true}

	db.Create(&menu)
	db.Create(&user)
	var users []User
	result := db.Where("name = ?", "Aka").First(&users)
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
	resultt := db.Where("name = ?", "salad").First(&menuse)
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
		if result := db.First(&dish, item.DishID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dish not found"})
			return
		}
		totalSum += float64(item.Quantity) * dish.Price
	}

	newOrder := Order{UserID: req.UserID, TotalSum: totalSum}
	if result := db.Create(&newOrder); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	for _, item := range req.CartItems {
		newOrderItem := OrderItem{OrderID: newOrder.ID, DishID: item.DishID, Quantity: item.Quantity}
		db.Create(&newOrderItem)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order_id": newOrder.ID})
}
