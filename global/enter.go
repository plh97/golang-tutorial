package global

import (
	"fmt"
	"golang-tutorial/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrate() {
	// if global.DB.Migrator().HasTable(&models.UserModel{}) {
	// 	global.DB.Migrator().DropTable(&models.UserModel{})
	// }
	err := DB.AutoMigrate(&models.UserModel{})
	if err != nil {
		log.Fatalln("Failed to migrate DB:", err)
	}

	// // 显式创建 UNIQUE 索引 (更可靠)
	// if !global.DB.Migrator().HasIndex(&models.UserModel{}, "Age") {
	// 	err = global.DB.Migrator().CreateIndex(&models.UserModel{}, "Age")
	// 	if err != nil {
	// 		log.Fatalln("Failed to create unique index on Name:", err)
	// 	}
	// }
}

func Connect() {
	// dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_db_new?parseTime=True"
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_db_new?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = db.Debug()
	if err != nil {
		fmt.Println(err)
		return
	}
	DB = db
}
