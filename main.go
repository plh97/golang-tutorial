package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    int
	Name  string
	Age   int
	Email string
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3305)/gorm_db_new"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		return
	}
	// Create a new User instance
	newUser := User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}

	// Create the record in the database
	err = db.Create(&newUser).Error
	if err != nil {
		log.Fatalf("create user failed: %s", err)
	}
	var userList []User
	db.Find(&userList)
	fmt.Println(userList)
}
