package main

import (
	"fmt"
	"golang-tutorial/global"
	"golang-tutorial/models"
	"log"
	"math/rand"
	"strconv"

	"gorm.io/gorm"
)

func query() {
	var userList models.UserModel
	global.DB.Unscoped().First(&userList, 11)
	fmt.Println(userList)
}

func create() {
	newUser := models.UserModel{
		Name: "John Doe" + strconv.Itoa(rand.Intn(10000)),
		// Email: "alice@example.com",
		Age: rand.Intn(50) + 1,
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
	global.DB.Unscoped().Delete(&models.UserModel{}, 14)
	// global.DB.Find(&models.UserModel{}, 11)
}

func deeplyQuery() {
	var userList []models.UserModel
	// query := global.DB.Where(models.UserModel{
	// 	Name: "2134",
	// 	Age:  0,
	// })
	// query := global.DB.Where("age > 18 and name = ?", "John Doe")
	// global.DB = global.DB.Where("name = ?", "John Doe")
	global.DB.Order("age asc").Not("age > 2").Or("name = ?", "John1 Doe").Find(&userList)
	fmt.Println(userList)
}

func scan() {
	// type User struct {
	// 	Label string `gorm:"column:name"`
	// 	Value int `gorm:"column:age"`
	// }
	// var nameList []User
	// global.DB.Model(models.UserModel{}).Scan(&nameList)
	// fmt.Println(nameList)

	var ageList []int
	global.DB.Model(models.UserModel{}).Distinct("age").Pluck("age", &ageList)
	fmt.Println(ageList)
}

func pagination() {
	limit := 2
	page := 0
	var userList []models.UserModel
	global.DB.Limit(limit).Offset((page - 1) * limit).Find(&userList)
	fmt.Println(userList)
}

func scope() {
	var userList []models.UserModel
	global.DB.Scopes(Age18).Find(userList)
	fmt.Println(userList)
}

func Age18(tx *gorm.DB) *gorm.DB {
	return tx.Where("age > 18")
}

func main() {
	global.Connect()
	global.Migrate()
	// create()
	// query()
	// update()
	// query()
	// delete()
	// deeplyQuery()
	// scan()
	// pagination()
	scope()
}
