package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryName string    `json:"categoryName" validate:"required,min=3"`
	CreatedAt    time.Time `json:"createdTime" gorm:"default:now()"`
	UpdatedAt    time.Time `json:"updatedTime" gorm:"default:now()"`
}

type CreateCategoryDTO struct {
	CategoryName string `json:"categoryName" validate:"required,min=3"`
}

type UpdateCategoryDTO struct {
	CategoryName string `json:"categoryName" validate:"required,min=3"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func MigrateCategory(db *gorm.DB) error {
	err := db.AutoMigrate(&Category{})
	return err
}
