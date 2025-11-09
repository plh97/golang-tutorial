package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        int64  `gorm:"primaryKey"`              // primary key
	Name      string `gorm:"size:25;not null,unique"` // name, unique, cannot be null
	Age       int    `gorm:"default 18"`
	Email     string
	CreatedAt time.Time //created by default
	DeletedAt gorm.DeletedAt
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
