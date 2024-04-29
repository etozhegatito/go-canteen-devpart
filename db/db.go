package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func ConnectDatabase() {
	dbs := "host=localhost user=postgres dbname=postgres password=mysecretpassword port=5432"
	db, err := gorm.Open(postgres.Open(dbs), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	} else {
		log.Println("Successfully connected to database")
	}

	db.AutoMigrate(&User{}, &Menu{}, OrderDetail{}, Order{})
	dateStr := "1990-12-31"

	// Преобразование строки в time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	user := User{Name: "Aka", Surname: "Aka", Password: "123", Phone: "111", Email: "profaka", DateOfBirth: date, IsAdmin: true}
	menu := Menu{Name: "salad", Price: 50, Description: "prosto", Weight: 11, InStock: true}
	orderDetail := OrderDetail{MenuID: 1, Count: 3, Menu: menu}
	db.Create(&menu)
	db.Create(&user)
	db.Create(&orderDetail)
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

	var menuse []Menu
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
