package db

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-canteen-devpart/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

// Для создание новго заказа
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
	var req models.OrderRequest
	if err := requests.ShouldBindJSON(&req); err != nil {
		requests.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Считоваем общую цену в корзине
	var totalSum float64
	for _, item := range req.CartItems {
		var dish models.Dish
		if result := database.First(&dish, item.DishID); result.Error != nil {
			requests.JSON(http.StatusBadRequest, gin.H{"error": "Dish not found"})
			return
		}
		totalSum += float64(item.Quantity) * dish.Price
	}

	//Создаем новый заказ в базе данных
	newOrder := models.Order{UserID: userID.(uint), TotalSum: totalSum}
	if result := database.Create(&newOrder); result.Error != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//Создаем OrderItem где указываем АЙДИ заказа и блюдо, количество
	for _, item := range req.CartItems {
		newOrderItem := models.OrderItem{OrderID: newOrder.ID, DishID: item.DishID, Quantity: item.Quantity}
		database.Create(&newOrderItem)
	}
	requests.JSON(http.StatusOK, gin.H{"message": "Харош, ты сделал заказ!", "order_id": newOrder.ID})
}

// Для создание нового блюда
func CreateDish(requests *gin.Context) {
	var dish models.Dish

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
	var dishes []models.Dish

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

	var dish models.Dish
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
	if result := database.Delete(&models.Dish{}, dishID); result.Error != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	requests.JSON(http.StatusOK, gin.H{"message": "Dish deleted successfully"})
}

func CreateUser(user models.User, requests *gin.Context) {
	if err := database.Create(&user).Error; err != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}
}

func CheckUser(creds models.Credentials, requests *gin.Context) {
	var user models.User
	log.Println("Checking user with email:", creds.Email)

	// Check if the user exists based on the provided email
	if result := database.Where("email = ?", creds.Email).First(&user); result.Error != nil {
		log.Println("User not found:", result.Error)
		requests.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Validate the user's password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Println("Password comparison failed:", err)
		requests.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Save the user session after successful authentication
	session := sessions.Default(requests)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		log.Println("Session save failed:", err)
		requests.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	log.Println("User authenticated successfully, session established")

	if user.IsAdmin {
		requests.Redirect(http.StatusFound, "/adminPage")
	} else {
		requests.Redirect(http.StatusFound, "/")
	}
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
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
func FetchOrdersWithDetails() ([]models.Order, error) {
	var orders []models.Order

	// Используем Preload для загрузки всех связанных данных
	if err := database.Preload("User").
		Preload("Items.Dish"). // Загрузка всех блюд через элементы заказа
		Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func OrdersHandler(c *gin.Context) {
	orders, err := FetchOrdersWithDetails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении заказов"})
		return
	}

	// Рендерим HTML-шаблон с полученными заказами
	c.HTML(http.StatusOK, "analystics.html", gin.H{
		"Orders": orders,
	})
}
