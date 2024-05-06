package db

import (
	"go-canteen-devpart/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
	database.AutoMigrate(&models.User{}, &models.Dish{}, &models.Order{}, &models.OrderItem{})

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
