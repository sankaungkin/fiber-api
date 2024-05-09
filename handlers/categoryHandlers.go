package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/dto"
	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/gorm"
	// ut "github.com/go-playground/universal-translator"
	// "github.com/go-playground/validator/v10"
)

// GetCategories godoc
//
//	@Summary		Fetch all Categories
//	@Description	Fetch all Categories
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Success		200				{array}		models.Category
//	@Failure		400				{object}	httputil.HttpError400
//	@Failure		401				{object}	httputil.HttpError401
//	@Failure		500				{object}	httputil.HttpError500
//	@Router			/api/category	[get]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func GetCategories(c *fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	println(user)
	claims := user.Claims.(jwt.MapClaims)
	// id := claims["id"].(string)
	id := claims["id"].(float64)
	c.SendString(fmt.Sprintf(`Hello user with id: %f`, id))

	db := database.DB

	categories := []models.Category{}

	db.Model(&models.Category{}).Order("ID asc").Limit(100).Find(&categories)

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

// UpdateCategory godoc
//
//	@Summary		Update individual category
//	@Description	Update individual category
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string							true	"category Id"
//	@Param			category			body		dto.UpdateCategoryRequestDTO	true	"Category Data"
//	@Success		200					{object}	models.Category
//	@Failure		400					{object}	httputil.HttpError400
//	@Failure		401					{object}	httputil.HttpError401
//	@Failure		500					{object}	httputil.HttpError500
//	@Router			/api/category/{id}	[put]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func UpdateCategory(c *fiber.Ctx) error {

	db := database.DB

	json := new(dto.UpdateCategoryRequestDTO)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON format",
		})
	}

	id := c.Params("id")

	var category models.Category
	result := db.First(&category, "id = ?", id)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "FAIL",
				"message": "No data",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status": "FAIL", "message": err.Error(),
		})
	}
	if json.CategoryName != "" {
		category.CategoryName = json.CategoryName
	}
	db.Save(&category)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Update successfully",
	})

}

// GetCategoryById godoc
//
//	@Summary		Fetch individual category by Id
//	@Description	Fetch individual category by Id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string	true	"category Id"
//	@Success		200					{object}	models.Category
//	@Failure		400					{object}	httputil.HttpError400
//	@Failure		401					{object}	httputil.HttpError401
//	@Failure		500					{object}	httputil.HttpError500
//	@Router			/api/category/{id}	[get]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func GetCategoryById(c *fiber.Ctx) error {
	db := database.DB

	id := c.Params("id")

	var category models.Category
	result := db.First(&category, "id = ?", id)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "FAIL",
				"message": "Record not found",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status": "FAIL", "message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Record found",
		"data":    category,
	})

}

// CreateCategory 	godoc
//
//	@Summary		Create new category based on parameters
//	@Description	Create new category based on parameters
//	@Tags			Categories
//	@Accept			json
//	@Param			category	body		dto.CreateCategoryRequestDTO	true	"Category Data"
//	@Success		200			{object}	models.Category
//	@Failure		400			{object}	httputil.HttpError400
//	@Failure		401			{object}	httputil.HttpError401
//	@Failure		500			{object}	httputil.HttpError500
//	@Failure		401			{object}	httputil.HttpError401
//	@Router			/api/category [post]
//
//	@Security		ApiKeyAuth
//
//	@param			Authorization	header	string	true	"Authorization"
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func CreateCategory(c *fiber.Ctx) error {

	db := database.DB

	input := new(dto.CreateCategoryRequestDTO)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid JSON format",
		})
	}

	category := models.Category{
		CategoryName: input.CategoryName,
	}

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

// DeleteCategory godoc
//
//	@Summary		Delete individual category
//	@Description	Delete individual category
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string	true	"category Id"
//	@Success		200					{object}	models.Category
//	@Failure		400					{object}	httputil.HttpError400
//	@Failure		401					{object}	httputil.HttpError401
//	@Failure		500					{object}	httputil.HttpError500
//	@Router			/api/category/{id}	[delete]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func DeleteCategory(c *fiber.Ctx) error {
	db := database.DB

	id := c.Params("id")

	var category models.Category

	result := db.First(&category, "id = ?", id)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "FAIL",
				"message": "Record not found",
			})
		}

		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  "ERROR",
			"message": err.Error(),
		})
	}
	db.Delete(&category)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Delete successfully",
	})

}
