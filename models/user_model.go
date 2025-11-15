package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID              int
	Name            string           `gorm:"size:16;not null;unique"` // name, unique, cannot be null
	Age             int              `gorm:"default:18"`
	UserDetailModel *UserDetailModel `gorm:"foreignKey:UserID"`
	CreatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}

type UserDetailModel struct {
	ID        int
	UserID    int        `gorm:"unique"`
	UserModel *UserModel `gorm:"foreignKey:UserID"`
	Email     string     `gorm:"size:32"`
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	fmt.Println("创建钩子函数")
	return nil
}

func (u *UserModel) BeforeUpdate(tx *gorm.DB) error {
	fmt.Println("update钩子函数")
	return nil
}

func (u *UserModel) BeforeDelete(tx *gorm.DB) error {
	fmt.Println("delete钩子函数")
	return nil
}
