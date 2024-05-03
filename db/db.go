package db

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// делаем базу данных глобальной переменной чтобы все функций имели доступ
var database *gorm.DB
var err error

func ConnectDatabase() {
	//Данные для входа в базу данных
	dbs := "host=db user=postgres dbname=postgres password=mysecretpassword port=5432"
	database, err = gorm.Open(postgres.Open(dbs), &gorm.Config{})

	//Проверяем валидность данных для подключение DataBase
	if err != nil {
		log.Fatal("Не получилось подключиться, Данные хуйня, исправь", err)
	} else {
		log.Println("Красавчик, база данных подключена!")
	}

	//Делаем миграцию чтобы изменение сохранились автоматом.
	database.AutoMigrate(&User{}, &Dish{}, &Order{}, &OrderItem{})

	//Для ручного ввода данных, если надо
	//user := User{Name: "Aqa", Surname: "Aqa", Password: "123", Phone: "111", Email: "pro", Age: 18, IsAdmin: true}
	//menu := Dish{Name: "salad", Price: 50, Description: "prosto", Weight: 11, InStock: true}
	//database.Create(&menu)
	//database.Create(&user)

	//Ручная проверка юзера, если надо
	//var users []User
	//result := database.Where("name = ?", "Aka").First(&users)
	//if result.Error != nil {
	//	if result.Error == gorm.ErrRecordNotFound {
	//		fmt.Println("No user found with the name Alice.")
	//	} else {
	//		log.Fatal("Error searching for user:", result.Error)
	//	}
	//} else {
	//	fmt.Printf("User found:")
	//}

	//Ручная проверка меню, если надо
	//var menuse []Dish
	//resultt := database.Where("name = ?", "salad").First(&menuse)
	//if resultt.Error != nil {
	//	if resultt.Error == gorm.ErrRecordNotFound {
	//		fmt.Println("No user found with the name Alice.")
	//	} else {
	//		log.Fatal("Error searching for user:", resultt.Error)
	//	}
	//} else {
	//	fmt.Printf("User found: %+v\n", menuse)
	//}
}

// Для создание заказа
func CreateOrder(requests *gin.Context) {
	//Получаем данные о сессий, нам нужно знать кто делает заказ
	session := sessions.Default(requests)
	userID := session.Get("user_id")

	//Вдруг не аутифицирован
	if userID == nil {
		requests.JSON(http.StatusUnauthorized, gin.H{"error": "Ты забыл зайти, дружок"})
		return
	}

	//Временный Struct для создание шаблона заказа
	var req OrderRequest
	if err := requests.ShouldBindJSON(&req); err != nil {
		requests.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Считоваем общую цену в корзине
	var totalSum float64
	for _, item := range req.CartItems {
		var dish Dish
		if result := database.First(&dish, item.DishID); result.Error != nil {
			requests.JSON(http.StatusBadRequest, gin.H{"error": "Dish not found"})
			return
		}
		totalSum += float64(item.Quantity) * dish.Price
	}

	//Создаем новый заказ в базе данных
	newOrder := Order{UserID: userID.(uint), TotalSum: totalSum}
	if result := database.Create(&newOrder); result.Error != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//Создаем OrderItem где указываем АЙДИ заказа и блюдо, количество
	for _, item := range req.CartItems {
		newOrderItem := OrderItem{OrderID: newOrder.ID, DishID: item.DishID, Quantity: item.Quantity}
		database.Create(&newOrderItem)
	}
	requests.JSON(http.StatusOK, gin.H{"message": "Харош, ты сделал заказ!", "order_id": newOrder.ID})
}

// Для создание нового блюда
func CreateDish(requests *gin.Context) {
	var dish Dish

	//Проверяем валидность данных с фронта
	if err := requests.ShouldBindJSON(&dish); err != nil {
		requests.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Создаем новое блюдо на базе данных и обработаем ошибку если будет
	if err := database.Create(&dish).Error; err != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Успех!
	requests.JSON(http.StatusCreated, dish)
}

// Возвращает весь список блюд в JSON формате
func GetDishes(requests *gin.Context) {
	var dishes []Dish

	//Ищем все блюда в базе данных и при форс мажоре узнаем из за чего ошибка
	if result := database.Find(&dishes); result.Error != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//Красавчик!
	requests.JSON(http.StatusOK, dishes)
}

// Обновление/редактирование существующего блюдо
func UpdateDish(requests *gin.Context) {
	//Получаем Айди блюдо
	id, err := strconv.Atoi(requests.Param("id"))
	if err != nil {
		requests.JSON(http.StatusBadRequest, gin.H{"error": "Неверный Айди блюдо дружок"})
		return
	}

	//
	var dish Dish
	if err := database.First(&dish, id).Error; err != nil {
		requests.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}

	if err := requests.ShouldBindJSON(&dish); err != nil {
		requests.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.Save(&dish)
	requests.JSON(http.StatusOK, dish)
}

func DeleteDish(requests *gin.Context) {
	dishID := requests.Param("id")
	if result := database.Delete(&Dish{}, dishID); result.Error != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	requests.JSON(http.StatusOK, gin.H{"message": "Dish deleted successfully"})
}

func CreateUser(user User, requests *gin.Context) {
	if err := database.Create(&user).Error; err != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}
}

func CheckUser(creds Credentials, requests *gin.Context) {
	var user User
	log.Println("Checking user with email:", creds.Email)

	if result := database.Where("email = ?", creds.Email).First(&user); result.Error != nil {
		log.Println("User not found:", result.Error)
		requests.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Println("Password comparison failed:", err)
		requests.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	session := sessions.Default(requests)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		log.Println("Session save failed:", err)
		requests.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	log.Println("User authenticated successfully, session established")

	//// Проверка, является ли пользователь администратором
	//if user.IsAdmin {
	//	token, err := auth.GenerateToken(user.ID, user.IsAdmin)
	//	if err != nil {
	//		requests.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	//		return
	//	}
	//	requests.JSON(http.StatusOK, gin.H{"token": token})
	//} else {
	//	requests.JSON(http.StatusUnauthorized, gin.H{"error": "Not admin"})
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
func GetAnalytics(requests *gin.Context) {
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
	requests.JSON(http.StatusOK, data)
}
