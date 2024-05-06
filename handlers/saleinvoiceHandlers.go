package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/gorm"
)

func CreateSaleInvoice(c *fiber.Ctx) error {

	type SaleInvoiceRequest struct {
		ID          string              `gorm:"primaryKey" json:"id"`
		CustomerId  uint                `json:"customerId"`
		SaleDetails []models.SaleDetail `gorm:"foreignKey:SaleId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"saleDetails"`
		Discount    int                 `json:"discount"`
		Total       int                 `json:"total"`
		GrandTotal  int                 `json:"grandTotal"`
		Remark      string              `json:"remark"`
		SaleDate    string              `jsong:"saleDate"`
	}

	db := database.DB

	input := new(SaleInvoiceRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid JSON format",
		})
	}

	newSale := models.Sale{}
	errors := models.ValidateStruct(newSale)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "operation failed",
		})
	}

	newSale.CustomerId = input.CustomerId
	newSale.Discount = input.Discount
	newSale.GrandTotal = input.GrandTotal
	newSale.ID = input.ID
	newSale.Remark = input.Remark
	newSale.SaleDate = input.SaleDate
	newSale.SaleDetails = input.SaleDetails
	newSale.Total = input.Total

	fmt.Println("NewSaleData : ", newSale)

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&newSale).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range newSale.SaleDetails {

		// decrease product qtyonhand
		var product models.Product
		result := tx.First(&product, "id = ?", newSale.SaleDetails[i].ProductId)
		if err := result.Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				tx.Rollback()
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status":  "FAIL",
					"message": "No data",
					"data":    err.Error(),
				})
			}
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"status":  "FAIL",
				"message": err.Error(),
			})
		}
		product.QtyOnHand -= int(newSale.SaleDetails[i].Qty)
		tx.Save(&product)

		// create inventory record
		newInventory := models.Inventory{
			InQty:     0,
			OutQty:    uint(newSale.SaleDetails[i].Qty),
			ProductId: newSale.SaleDetails[i].ProductId,
			Remark:    "SaleID:" + newSale.ID + ", line items id:" + strconv.Itoa(int(newSale.SaleDetails[i].ID)) + ", decrease quantity: " + strconv.Itoa(newSale.SaleDetails[i].Qty),
		}
		tx.Save(&newInventory)

	}
	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": newSale,
	})

}
