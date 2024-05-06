package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreatePurchaseInvoice(c *fiber.Ctx) error {

	type PurchaseInvoiceRequest struct {
		ID              string                  `gorm:"primaryKey" json:"id"`
		SupplierId      uint                    `json:"supplierId"`
		PurchaseDetails []models.PurchaseDetail `gorm:"foreignKey:PurchaseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"purchaseDetails"`
		Discount        int                     `json:"discount"`
		Total           int                     `json:"total"`
		GrandTotal      int                     `json:"grandTotal"`
		Remark          string                  `json:"remark"`
		PurchaseDate    string                  `jsong:"purchaseDate"`
	}

	db := database.DB

	input := new(PurchaseInvoiceRequest)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid JSON format",
		})
	}

	newPurchase := models.Purchase{
		ID:              input.ID,
		SupplierId:      input.SupplierId,
		PurchaseDetails: input.PurchaseDetails,
		Discount:        input.Discount,
		Total:           input.Total,
		GrandTotal:      input.GrandTotal,
		Remark:          input.Remark,
		PurchaseDate:    input.PurchaseDate,
	}

	errors := models.ValidateStruct(newPurchase)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "operation failed",
		})
	}

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&newPurchase).Error; err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("newPruchase : ", newPurchase)
	for i := range newPurchase.PurchaseDetails {

		//increase product qtyonhand
		var product models.Product
		result := tx.First(&product, "id = ?", newPurchase.PurchaseDetails[i].ProductId)
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
		product.QtyOnHand += newPurchase.PurchaseDetails[i].Qty
		tx.Save(&product)

		//create inventory records
		newInventory := models.Inventory{
			InQty:     uint(newPurchase.PurchaseDetails[i].Qty),
			OutQty:    0,
			ProductId: newPurchase.PurchaseDetails[i].ProductId,
			Remark:    "PurchaseID:" + newPurchase.ID + ", line items id: " + strconv.Itoa(int(newPurchase.PurchaseDetails[i].ID)) + ", increase quantity: " + strconv.Itoa(newPurchase.PurchaseDetails[i].Qty),
		}
		tx.Save(&newInventory)

		// create item transaction records
		newItemTransaction := models.ItemTransaction{
			ProductId:   newPurchase.PurchaseDetails[i].ProductId,
			ReferenceNo: newPurchase.ID + "-" + strconv.Itoa(int(newPurchase.PurchaseDetails[i].ID)),
			InQty:       newPurchase.PurchaseDetails[i].Qty,
			OutQty:      0,
			TranType:    "DR",
			Remark:      "PurchaseID:" + newPurchase.ID + ", line items id: " + strconv.Itoa(int(newPurchase.PurchaseDetails[i].ID)) + ", increase quantity: " + strconv.Itoa(newPurchase.PurchaseDetails[i].Qty),
		}
		tx.Save(&newItemTransaction)
	}
	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": newPurchase,
	})
}

func GetPurchases(c *fiber.Ctx) error {
	db := database.DB

	purchases := []models.Purchase{}

	// db.Model(&models.Purchase{}).Order("ID asc").Preload("PurchaseDetails").Find(&purchases)
	db.Preload(clause.Associations).Find(&purchases)

	if len(purchases) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "NO RECORD",
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": len(purchases),
		"data":    purchases,
	})

}
