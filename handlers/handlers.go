package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/models"
)

// type Repository struct{
// 	DB *gorm.DB
// }

type Category struct {
	
	CategoryName	string	`json:"categoryName"`

}

func GetCategories(c *fiber.Ctx) error {
	
	db := database.DB

	categories := []models.Category{}

 err := db.Model(&models.Category{}).Order("id asc").Limit(100).Find(&categories)

if len(categories) == 0 {
	return c.Status(200).JSON(fiber.Map{
		"status": "SUCCESS with no data", 
		"message": "No Data", 
		"data": nil,
	})
}

	if err != nil {
		return err.Error
	}

	return c.Status(200).JSON(fiber.Map{
		"code": 200, 
		"message": "SUCCESS", 
		"data": categories,
	})

}

func CreateCategory(c *fiber.Ctx) error {

	db :=database.DB

	category := Category{}

	// validate the CREATE CATEGORY DTO
	err := c.BodyParser(&category)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{ "message":"request failed"})
		return err
	}

	err = db.Create(&category).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{ "message": "could not create new category"})
		return err
	}
	 
	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"category has been created successfully"})
	return nil
	}
