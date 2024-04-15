package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID				uint		`gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryName	string		`json:"categoryName"`
	CreatedAt		time.Time	`json:"createdTime" gorm:"default:now()"`
	UpdatedAt		time.Time 	`json:"updatedTime" gorm:"default:now()"`
}

type NewCategory struct {
	gorm.Model
	CategoryName string 	`gorm:"uniqueIndex"`
}

func MigrateCategory(db *gorm.DB) error {
	err := db.AutoMigrate(&Category{}, &NewCategory{})
	return err
}