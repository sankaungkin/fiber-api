package dto

import "github.com/sankaungkin/fiber-api/models"

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email" example:"user@email.com"`
	Password string `json:"password" validate:"required" example:"pass1234"`
}

type LoginResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type CreateUserRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
	Role     string `json:"role"`
}

type UpdateCategoryRequestDTO struct {
	CategoryName string `json:"categoryName"`
}

type CreateCategoryRequestDTO struct {
	CategoryName string `json:"categoryName"`
}

type CreateProductRequstDTO struct {
	ID              string `gorm:"primaryKey" json:"id"`
	ProductName     string `json:"productName" validate:"required,min=3"`
	CategoryId      uint   `json:"categoryId" validate:"required"`
	Uom             string `json:"uom" validate:"required,min=3"`
	BuyPrice        int16  `josn:"buyPrice" validate:"required,min=1"`
	SellPriceLevel1 int16  `josn:"sellPricelvl1" validate:"required,min=1"`
	SellPriceLevel2 int16  `josn:"sellPricelvl2" validate:"required,min=1"`
	ReorderLvl      uint   `json:"reorderlvl" gorm:"default:1" validate:"required,min=1"`
	QtyOnHand       int    `json:"qtyOhHand" validate:"required"`
	BrandName       string `json:"brand"`
	IsActive        bool   `json:"isActive" gorm:"default:true"`
}

type UpdateProductRequstDTO struct {
	ProductName     string `json:"productName" validate:"required,min=3"`
	CategoryId      uint   `json:"categoryId" validate:"required"`
	Uom             string `json:"uom" validate:"required,min=2"`
	BuyPrice        int16  `josn:"buyPrice" validate:"required,min=1"`
	SellPriceLevel1 int16  `josn:"sellPricelvl1" validate:"required,min=1"`
	SellPriceLevel2 int16  `josn:"sellPricelvl2" validate:"required,min=1"`
	ReorderLvl      uint   `json:"reorderlvl" gorm:"default:1" validate:"required,min=1"`
	// QtyOnHand       int    `json:"qtyOhHand" validate:"required"`
	BrandName string `json:"brand"`
	IsActive  bool   `json:"isActive" gorm:"default:true"`
}

type SaleInvoiceRequestDTO struct {
	ID          string              `gorm:"primaryKey" json:"id"`
	CustomerId  uint                `json:"customerId"`
	SaleDetails []models.SaleDetail `gorm:"foreignKey:SaleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"saleDetails"`
	Discount    int                 `json:"discount"`
	Total       int                 `json:"total"`
	GrandTotal  int                 `json:"grandTotal"`
	Remark      string              `json:"remark"`
	SaleDate    string              `jsong:"saleDate"`
}

type PurchaseInvoiceRequestDTO struct {
	ID              string                  `gorm:"primaryKey" json:"id"`
	SupplierId      uint                    `json:"supplierId"`
	PurchaseDetails []models.PurchaseDetail `gorm:"foreignKey:PurchaseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"purchaseDetails"`
	Discount        int                     `json:"discount"`
	Total           int                     `json:"total"`
	GrandTotal      int                     `json:"grandTotal"`
	Remark          string                  `json:"remark"`
	PurchaseDate    string                  `jsong:"purchaseDate"`
}
