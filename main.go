package main

import (
	"fmt"
	"golang-tutorial/global"
	"golang-tutorial/models"
	"log"
	"math/rand"
)

func query() {
	var userList models.UserModel
	global.DB.Unscoped().First(&userList, 11)
	fmt.Println(userList)
}

func create() {
	newUser := models.UserModel{
		Name:  "John Doe",
		Email: "alice@example.com",
		Age:   rand.Intn(50) + 1,
	}

	//Create the record in the database
	err := global.DB.Create(&newUser).Error
	if err != nil {
		log.Fatalf("create user failed: %s", err)
	}
}

func update() {
	newUser := models.UserModel{
		// Name:  "John Doe112",
		// Email: "alice@example.com",
		Age: 0,
		// ID:        5,
		// CreatedAt: time.Now(),
	}

	//Create the record in the database
	err := global.DB.Model(&models.UserModel{ID: 1}).Select("age").Updates(newUser).Error
	if err != nil {
		log.Fatalf("create user failed: %s", err)
	}
}
func delete() {
	global.DB.Unscoped().Delete(&models.UserModel{}, 11)
	// global.DB.Find(&models.UserModel{}, 11)
}

func deeplyQuery() {
	var userList []models.UserModel
	global.DB.Where(models.UserModel{
		Name: "2134",
		Age:  0,
	}).First(&userList)
	fmt.Println(userList)
}

func main() {
	global.Connect()
	global.Migrate()
	create()
	query()
	update()
	// query()
	delete()
	deeplyQuery()
}
