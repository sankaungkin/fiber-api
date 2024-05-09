package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/dto"
	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreatePurchaseInvoice 	godoc
//
//	@Summary		Create new Purchase Invoice based on parameters
//	@Description	Create new Purchase Invoice based on parameters
//	@Tags			Purchases
//	@Accept			json
//	@Param			purchase	invoice		body	dto.PurchaseInvoiceRequestDTO	true	"purchase invoice Data"
//	@Success		200			{object}	models.Purchase
//	@Failure		400			{object}	httputil.HttpError400
//	@Failure		401			{object}	httputil.HttpError401
//	@Failure		500			{object}	httputil.HttpError500
//	@Failure		401			{object}	httputil.HttpError401
//	@Router			/api/purchase [post]
//
//	@Security		ApiKeyAuth
//
//	@param			Authorization	header	string	true	"Authorization"
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func CreatePurchaseInvoice(c *fiber.Ctx) error {

	db := database.DB

	input := new(dto.PurchaseInvoiceRequestDTO)

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

// GetPurchases godoc
//
//	@Summary		Fetch all purchase invoices
//	@Description	Fetch all purchase invoices
//	@Tags			Purchases
//	@Accept			json
//	@Produce		json
//	@Success		200				{array}		models.Purchase
//	@Failure		400				{object}	httputil.HttpError400
//	@Failure		401				{object}	httputil.HttpError401
//	@Failure		500				{object}	httputil.HttpError500
//	@Router			/api/purchase	[get]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
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

// GetPurchaseById godoc
//
//	@Summary		Fetch individual purchase invoice by Id
//	@Description	Fetch individual purchase invoice by Id
//	@Tags			Purchases
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string	true	"purchase Id"
//	@Success		200					{object}	models.Purchase
//	@Failure		400					{object}	httputil.HttpError400
//	@Failure		401					{object}	httputil.HttpError401
//	@Failure		500					{object}	httputil.HttpError500
//	@Router			/api/purchase/{id}	[get]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func GetPurchaseById(c *fiber.Ctx) error {

	db := database.DB

	id := strings.ToUpper(c.Params("id"))

	var purchase models.Purchase

	result := db.First(&purchase, "id = ?", id)

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
		"data":    purchase,
	})

}
