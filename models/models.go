package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Category struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryName string    `json:"categoryName" validate:"required,min=3"`
	Products     []Product `gorm:"foreignKey:CategoryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt    time.Time `json:"createdTime" gorm:"default:now()"`
	UpdatedAt    time.Time `json:"updatedTime" gorm:"default:now()"`
}

type Product struct {
	gorm.Model
	ID              string      `gorm:"primaryKey" json:"id"`
	ProductName     string      `json:"productName" validate:"required,min=3"`
	CategoryId      uint        `json:"categoryId"`
	Inventories     []Inventory `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Uom             string      `json:"uom" validate:"required,min=3"`
	BuyPrice        int16       `josn:"buyPrice" validate:"required,min=1"`
	SellPriceLevel1 int16       `josn:"sellPricelvl1" validate:"required,min=1"`
	SellPriceLevel2 int16       `josn:"sellPricelvl2" validate:"required,min=1"`
	ReorderLvl      uint        `json:"reorderlvl" gorm:"default:1" validate:"required,min=1"`
	QtyOnHand       int         `json:"qtyOhHand" validate:"required"`
	BrandName       string      `json:"brand"`
	IsActive        bool        `json:"isActive" gorm:"default:true"`
	CreatedAt       int64       `gorm:"autoCreateTime" json:"-"`
	UpdatedAt       int64       `gorm:"autoUpdateTime:milli" json:"-"`
}

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string `gorm:"uniqueIndex;" json:"email" validate:"required,email"`
	UserName  string `json:"username" validate:"required,min=3"`
	Password  string `json:"password" validate:"required,min=3"`
	IsAdmin   bool   `json:"isAdmin" validate:"required"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"-"`
}

type Customer struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey:autoIncrement" json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Sales     []Sale `gorm:"foreignKey:CustomerId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"-"`
}

type Supplier struct {
	gorm.Model
	ID        uint       `gorm:"primaryKey:autoIncrement" json:"id"`
	Name      string     `json:"name"`
	Address   string     `json:"address"`
	Phone     string     `json:"phone"`
	Purchases []Purchase `gorm:"foreignKey:SupplierId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt int64      `gorm:"autoUpdateTime:milli" json:"-"`
}

type Purchase struct {
	gorm.Model
	ID              string           `gorm:"primaryKey" json:"id"`
	SupplierId      uint             `json:"-"`
	PurchaseDetails []PurchaseDetail `gorm:"foreignKey:PurchaseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Discount        int              `json:"discount"`
	Total           int              `json:"total"`
	GrandTotal      int              `json:"grandTotal"`
	Remark          string           `json:"remark"`
	PurchaseDate    string           `jsong:"purchaseDate"`
	CreatedAt       int64            `gorm:"autoCreateTime" json:"-"`
	UpdatedAt       int64            `gorm:"autoUpdateTime:milli" json:"-"`
}

type PurchaseDetail struct {
	gorm.Model
	ID          string `gorm:"primaryKey" json:"id"`
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	Qty         int    `json:"qty"`
	Price       int    `json:"price"`
	Total       int    `json:"total"`
	PurchaseId  string `json:"-"`
}

type Inventory struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey:autoIncrement" json:"id"`
	OutQty    uint   `json:"inQty"`
	InQty     uint   `json:"outQty"`
	ProductId string `json:"productId"`
	Remark    string `json:"remark"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"-"`
}

type Sale struct {
	gorm.Model
	ID           string       `gorm:"primaryKey" json:"id"`
	CustomerId   uint         `json:"-"`
	SaleDetails  []SaleDetail `gorm:"foreignKey:SaleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Discount     int          `json:"discount"`
	Total        int          `json:"total"`
	GrandTotal   int          `json:"grandTotal"`
	Remark       string       `json:"remark"`
	PurchaseDate string       `jsong:"purchaseDate"`
	CreatedAt    int64        `gorm:"autoCreateTime" json:"-"`
	UpdatedAt    int64        `gorm:"autoUpdateTime:milli" json:"-"`
}

type SaleDetail struct {
	gorm.Model
	ID          string `gorm:"primaryKey" json:"id"`
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	Qty         int    `json:"qty"`
	Price       int    `json:"price"`
	Total       int    `json:"total"`
	SaleId      string `json:"-"`
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
	err := db.AutoMigrate(&Category{}, &Product{}, &User{}, &Customer{}, &Supplier{}, &Sale{}, &SaleDetail{}, &Purchase{}, &PurchaseDetail{}, &User{})
	return err
}
