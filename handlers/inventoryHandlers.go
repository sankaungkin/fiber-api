package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/gorm"
	// ut "github.com/go-playground/universal-translator"
	// "github.com/go-playground/validator/v10"
)

func IncreaseInventory(c *fiber.Ctx) error {

	type InventoryRequest struct {
		InQty     uint   `json:"inQty" validate:"required, min=0"`
		OutQty    uint   `json:"outQty" validate:"required, min=0"`
		ProductId string `json:"productId" validate:"required"`
		Remark    string `json:"remark"`
	}
	db := database.DB

	input := new(InventoryRequest)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid JSON format",
		})
	}
	newInventory := models.Inventory{}

	errors := models.ValidateStruct(newInventory)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "operation failed",
		})
	}

	newInventory.InQty = input.InQty
	newInventory.OutQty = input.OutQty
	newInventory.ProductId = input.ProductId
	newInventory.Remark = input.Remark
	// start transaction

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	// create new record to inventory

	if err := tx.Create(&newInventory).Error; err != nil {
		tx.Rollback()
		return err
	}
	// increate QtyOnHand value for respected productId

	var product models.Product
	result := tx.First(&product, "id = ?", input.ProductId)
	// fmt.Printf("result: ", result)
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
	product.QtyOnHand += int(input.InQty)
	message := input.ProductId + " is increased by " + strconv.Itoa(int(input.InQty)) + "EACH"
	tx.Save(&product)

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": message,
	})
}

func DecreaseInventory(c *fiber.Ctx) error {

	type InventoryRequest struct {
		InQty     uint   `json:"inQty" validate:"required, min=0"`
		OutQty    uint   `json:"outQty" validate:"required, min=0"`
		ProductId string `json:"productId" validate:"required"`
		Remark    string `json:"remark"`
	}
	db := database.DB

	input := new(InventoryRequest)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid JSON format",
		})
	}
	newInventory := models.Inventory{}

	errors := models.ValidateStruct(newInventory)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "operation failed",
		})
	}

	newInventory.InQty = input.InQty
	newInventory.OutQty = input.OutQty
	newInventory.ProductId = input.ProductId
	newInventory.Remark = input.Remark
	// start transaction

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	// create new record to inventory

	if err := tx.Create(&newInventory).Error; err != nil {
		tx.Rollback()
		return err
	}
	// decrease QtyOnHand value for respected productId

	var product models.Product
	result := tx.First(&product, "id = ?", input.ProductId)
	// fmt.Printf("result: ", result)
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
	product.QtyOnHand -= int(input.OutQty)
	message := input.ProductId + " is decreased by " + strconv.Itoa(int(input.OutQty)) + "EACH"
	tx.Save(&product)

	tx.Commit()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": message,
	})
}
