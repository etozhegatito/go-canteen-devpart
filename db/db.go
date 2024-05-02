package db

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

var jwtKey = []byte("your_secret_key") // Используйте безопасный способ хранения ключа
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

	//user := User{Name: "Aqa", Surname: "Aqa", Password: "123", Phone: "111", Email: "pro", Age: 18, IsAdmin: true}
	//menu := Dish{Name: "salad", Price: 50, Description: "prosto", Weight: 11, InStock: true}

	//database.Create(&menu)
	//database.Create(&user)

	var users []User
	result := database.Where("name = ?", "Aka").First(&users)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("No user found with the name Alice.")
		} else {
			log.Fatal("Error searching for user:", result.Error)
		}
	} else {
		fmt.Printf("User found:")
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
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
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

	newOrder := Order{UserID: userID.(uint), TotalSum: totalSum}
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

func CreateDish(c *gin.Context) {
	var dish Dish
	if err := c.ShouldBindJSON(&dish); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.Create(&dish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dish)
}

func GetDishes(c *gin.Context) {
	var dishes []Dish
	if result := database.Find(&dishes); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, dishes)
}

func UpdateDish(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var dish Dish
	if err := database.First(&dish, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}

	if err := c.ShouldBindJSON(&dish); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Save(&dish)
	c.JSON(http.StatusOK, dish)
}

func DeleteDish(c *gin.Context) {
	dishID := c.Param("id")
	if result := database.Delete(&Dish{}, dishID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dish deleted successfully"})
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

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		log.Println("Session save failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	log.Println("User authenticated successfully, session established")

	//// Проверка, является ли пользователь администратором
	//if user.IsAdmin {
	//	token, err := auth.GenerateToken(user.ID, user.IsAdmin)
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	//		return
	//	}
	//	c.JSON(http.StatusOK, gin.H{"token": token})
	//} else {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Not admin"})
	//}
}

func GetUserByID(id uint) (*User, error) {
	var user User
	// Отправка запроса в базу данных и загрузка первой найденной записи в `user`
	result := database.Where("id = ?", id).First(&user)
	if result.Error != nil {
		// Возвращаем nil и ошибку, если пользователь не найден или произошла другая ошибка
		return nil, result.Error
	}
	// Возвращаем найденного пользователя и nil в качестве ошибки, если пользователь найден успешно
	return &user, nil
}

// Модуль db должен иметь функции для получения статистики

// В db.go
func GetAnalytics(c *gin.Context) {
	data := AnalyticsData{
		TotalOrders:  150,
		TotalRevenue: 12345.67,
		AverageBill:  82.30,
		TotalUsers:   75,
		MostPopularDish: DishStat{
			DishID:   1,
			DishName: "Spaghetti Carbonara",
			Count:    120,
		},
		LeastPopularDishes: []DishStat{
			{DishID: 2, DishName: "Brussels Sprouts", Count: 5},
			{DishID: 3, DishName: "Canned Tuna Salad", Count: 8},
			{DishID: 4, DishName: "Burnt Toast", Count: 10},
		},
	}

	// Возвращаем данные в формате JSON
	c.JSON(http.StatusOK, data)
}
