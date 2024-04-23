package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/models"
)

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
