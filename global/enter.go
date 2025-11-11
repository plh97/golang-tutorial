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
	if DB.Migrator().HasTable(&models.UserDetailModel{}) {
		DB.Migrator().DropTable(&models.UserDetailModel{})
	}
	err := DB.AutoMigrate(&models.UserModel{}, &models.UserDetailModel{})
	if err != nil {
		log.Fatalln("Failed to migrate DB:", err)
	}
}

func Connect() {
	// dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_db_new?parseTime=True"
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_db_new?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	db = db.Debug()
	if err != nil {
		fmt.Println(err)
		return
	}
	DB = db
}
