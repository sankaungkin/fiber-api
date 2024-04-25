package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/gorm"
)

func CreateProduct(c *fiber.Ctx) error {

	db := database.DB

	newProduct := models.Product{}

	err := c.BodyParser(&newProduct)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "request failed",
		})
		return err
	}

	errors := models.ValidateStruct(newProduct)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validations failed",
		})
	}

	err = db.Create(&newProduct).Error
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "operation failed",
		})
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SUCCESS",
		"data":    newProduct,
	})

}

func GetProducts(c *fiber.Ctx) error {
	db := database.DB

	products := []models.Product{}

	db.Model(&models.Product{}).Order("ID asc").Find(&products)

	if len(products) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "NO RECORD",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": len(products),
		"data":    products,
	})
}

func GetProductById(c *fiber.Ctx) error {
	db := database.DB

	id := c.Params("id")

	var product models.Product
	result := db.First(&product, "id = ?", id)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "FAIL",
				"message": "Record not found",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "FAIL",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Record found",
		"data":    product,
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	type UpdateProductRequest struct {
		ProductName     string `json:"productName" validate:"required,min=3"`
		CategoryId      uint   `json:"categoryId" validate:"required"`
		Uom             string `json:"uom" validate:"required,min=2"`
		BuyPrice        int16  `josn:"buyPrice" validate:"required,min=1"`
		SellPriceLevel1 int16  `josn:"sellPricelvl1" validate:"required,min=1"`
		SellPriceLevel2 int16  `josn:"sellPricelvl2" validate:"required,min=1"`
		ReorderLvl      uint   `json:"reorderlvl" gorm:"default:1" validate:"required,min=1"`
		QtyOnHand       int    `json:"qtyOhHand" validate:"required"`
		BrandName       string `json:"brand"`
		IsActive        bool   `json:"isActive" gorm:"default:true"`
	}

	db := database.DB

	input := new(UpdateProductRequest)
	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON format",
		})
	}

	id := c.Params("id")

	var product models.Product
	result := db.First(&product, "id = ?", id)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "FAIL",
				"message": "No data",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "FAIL",
			"message": err.Error(),
		})
	}

	if input.ProductName != "" || input.Uom != "" {
		product.BrandName = input.BrandName
		product.BuyPrice = input.BuyPrice
		product.CategoryId = input.CategoryId
		product.IsActive = input.IsActive
		product.ProductName = input.ProductName
		product.QtyOnHand = input.QtyOnHand
		product.ReorderLvl = input.ReorderLvl
		product.SellPriceLevel1 = input.SellPriceLevel1
		product.SellPriceLevel2 = input.SellPriceLevel2
		product.Uom = input.Uom
	}

	db.Save(&product)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Update successfully",
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	db := database.DB

	id := c.Params("id")

	var product models.Product

	result := db.First(&product, "id = ?", id)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "FAIL",
				"message": "Record not found",
			})
		}

		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"stats":   "ERROR",
			"message": err.Error(),
		})
	}

	db.Delete(&product)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Delete successfully",
	})
}
