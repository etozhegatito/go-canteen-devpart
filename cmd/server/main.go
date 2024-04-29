package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	//http.HandleFunc("/register", auth.RegisterHandler)
	//http.HandleFunc("/login", auth.LoginHandler)
	//
	//log.Println("Starting server on :8080...")
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}
	dbs := "host=localhost user=postgres dbname=postgres password=mysecretpassword port=5432"
	_, err := gorm.Open(postgres.Open(dbs), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
}
