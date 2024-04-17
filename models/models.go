package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	guuid "github.com/google/uuid"
)

type Category struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryName string    `json:"categoryName" validate:"required,min=3"`
	CreatedAt    time.Time `json:"createdTime" gorm:"default:now()"`
	UpdatedAt    time.Time `json:"updatedTime" gorm:"default:now()"`
}

type User struct {
	gorm.Model

	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	UserName  string    `json:"username" validate:"required,min=3"`
	Password  string    `json:"password" validate:"required,min=3"`
	Session   []Session `gorm:"foreignKey:UserRefer;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"-"`
	UpdatedAt int64     `gorm:"autoUpdateTime:milli" json:"-"`
}

type Session struct {
	Sessionid guuid.UUID `gorm:"primaryKey" json:"sessionid"`
	Expires   time.Time  `json:"-"`
	UserRefer uint       `json:"-"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-"`
}

type ErrorResponse struct {
	Field string                                 `json:"field"`
	Tag   string                                 `json:"tag"`
	Value string                                 `json:"value,omitempty"`
	Info  validator.ValidationErrorsTranslations `json:"info"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {

	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	var errors []*ErrorResponse
	err := validate.Struct(payload)

	if err != nil {

		errTran := err.(validator.ValidationErrors)
		fmt.Println(errTran.Translate(trans))
		info := errTran.Translate(trans)

		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.Info = info
			errors = append(errors, &element)
		}

	}
	return errors
}

func MigrateModels(db *gorm.DB) error {
	err := db.AutoMigrate(&Category{}, &User{})
	return err
}
