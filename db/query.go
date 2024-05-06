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
	//Для откладки
	if err := database.First(&dish, id).Error; err != nil {
		requests.JSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
		return
	}

	if err := requests.ShouldBindJSON(&dish); err != nil {
		requests.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Сохроняем изменение в базе данных
	database.Save(&dish)
	requests.JSON(http.StatusOK, dish)
}

// Для удаление блюдо из база данных
func DeleteDish(requests *gin.Context) {
	//получаем айди блюдо и удалаем блюдо
	dishID := requests.Param("id")
	if result := database.Delete(&models.Dish{}, dishID); result.Error != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//Успех!
	requests.JSON(http.StatusOK, gin.H{"message": "Dish deleted successfully"})
}

// Создаем нового пользователя
func CreateUser(user models.User, requests *gin.Context) {
	//При форс мажорах узнаем в чем ошибка
	if err := database.Create(&user).Error; err != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}
}

// Проверяем наличие юзера в базе данных
func CheckUser(creds models.Credentials, requests *gin.Context) {
	var user models.User
	//Для откладки
	log.Println("Checking user with email:", creds.Email)

	// Ищем юзера по его емайлу
	if result := database.Where("email = ?", creds.Email).First(&user); result.Error != nil {
		log.Println("User not found:", result.Error)
		requests.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	//Хэшируем пароль и проверяем его правильность
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		log.Println("Password comparison failed:", err)
		requests.JSON(http.StatusUnauthorized, gin.H{"error": "Неправильный пароль ты выбрал дружок!"})
		return
	}

	//После успешного захода сохроняем сессию
	session := sessions.Default(requests)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		log.Println("Сессия не сохранилась", err)
		requests.JSON(http.StatusInternalServerError, gin.H{"error": "Провал при сохранений сессий"})
		return
	}
	//Для проверки добавил логи
	log.Println("Юзер успешно авторизовался!")

	//Перенаправляем юзера
	//Если админ - в админску страницу. Если обычный юзер то в Dashboard
	if user.IsAdmin {
		requests.Redirect(http.StatusFound, "/adminPage")
	} else {
		requests.Redirect(http.StatusFound, "/")
	}
}

// Просто находим юзера по его айди
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	// Отправка запроса в базу данных и загрузка первой найденной записи в `user`
	result := database.Where("id = ?", id).First(&user)
	if result.Error != nil {
		// Возвращаем nil и ошибку, если пользователь не найден или произошла другая ошибка
		return nil, result.Error
	}
	// Возвращаем найденного пользователя и nil, если пользователь найден успешно
	return &user, nil
}

// Для админа, общая статистика такие как средний чек, общая сумма чеков
func ReportPage() ([]models.Order, float64, float64, error) {
	var orders []models.Order
	var totalSum float64
	var averageCheck float64

	//Ищем все все заказы в базе данных
	if err := database.Find(&orders).Error; err != nil {
		return nil, 0, 0, err
	}

	//высчитоваем общую сумму заказов
	for _, order := range orders {
		totalSum += order.TotalSum
	}

	//если заказ существует тогда сделаем рассчет среднего чека
	if len(orders) > 0 {
		averageCheck = totalSum / float64(len(orders))
	}

	return orders, totalSum, averageCheck, nil
}

// Для репорт страницы
func OrdersHandler(requests *gin.Context) {
	orders, totalSum, avgCheck, err := ReportPage()
	if err != nil {
		requests.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении заказов"})
		return
	}

	// Передаем заказы, общую сумму и средний чек в шаблон
	requests.HTML(http.StatusOK, "analystics.html", gin.H{
		"Orders":   orders,
		"TotalSum": totalSum,
		"AvgCheck": avgCheck,
	})
}
