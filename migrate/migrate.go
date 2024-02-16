package main

import (
	"echo-todo-app/db"
	"echo-todo-app/model"
	"log"
)

func main() {
	dbConn := db.NewDB()
	defer log.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.Task{})
}
