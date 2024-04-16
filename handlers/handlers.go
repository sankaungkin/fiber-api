package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/gorm"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Category struct {
	CategoryName string `json:"categoryName"`
}

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func GetCategories(c *fiber.Ctx) error {

	db := database.DB

	categories := []models.Category{}

	db.Model(&models.Category{}).Order("ID asc").Limit(100).Find(&categories)

	fmt.Println(len(categories))

	if len(categories) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    "404",
			"message": "NO RECORD",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": len(categories),
		"data":    categories,
	})

}

func GetCategory(c *fiber.Ctx) error {
	db := database.DB

	id := c.Params("id")

	var category models.Category
	result := db.First(&category, "id = ?", id)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "fail",
				"message": "No note with that Id exists",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status": "fail", "message": err.Error(),
		})
	}

	fmt.Println(&result)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Record has been found",
		"data":    category,
	})

}

func CreateCategory(c *fiber.Ctx) error {

	db := database.DB

	category := models.Category{}

	// newCategory := models.CreateCategoryDTO{}

	// validate the CREATE CATEGORY DTO
	err := c.BodyParser(&category)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	errors := models.ValidateStruct(category)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	err = db.Create(&category).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create new category"})
		return err
	}

	return c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status":  "SUCCESS",
			"message": "category has been created successfully",
			"data":    category,
		})

}
