package main

import (
	"go-canteen-devpart/db"
)

func main() {
	//http.HandleFunc("/register", auth.RegisterHandler)
	//http.HandleFunc("/login", auth.LoginHandler)
	//
	//log.Println("Starting server on :7070...")
	//if err := http.ListenAndServe(":7070", nil); err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}
	db.ConnectDatabase()

}
