package handlers

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	_ "github.com/sankaungkin/fiber-api/cmd/docs"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/dto"
	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/gorm"
)

// CreateProduct 	godoc
//
//	@Summary		Create new product based on parameters
//	@Description	Create new product based on parameters
//	@Tags			Products
//	@Accept			json
//	@Param			product	body		dto.CreateProductRequstDTO	true	"Product Data"
//	@Success		200		{object}	models.Product
//	@Failure		400		{object}	httputil.HttpError400
//	@Failure		401		{object}	httputil.HttpError401
//	@Failure		500		{object}	httputil.HttpError500
//	@Failure		401		{object}	httputil.HttpError401
//	@Router			/api/product [post]
//
//	@Security		ApiKeyAuth
//
//	@param			Authorization	header	string	true	"Authorization"
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func CreateProduct(c *fiber.Ctx) error {

	db := database.DB

	input := new(dto.CreateProductRequstDTO)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "Invalid JSON format",
		})
	}

	newProduct := models.Product{
		ID:              input.ID,
		ProductName:     input.ProductName,
		CategoryId:      input.CategoryId,
		Uom:             input.Uom,
		BuyPrice:        input.BuyPrice,
		SellPriceLevel1: input.SellPriceLevel1,
		SellPriceLevel2: input.SellPriceLevel2,
		ReorderLvl:      input.ReorderLvl,
		QtyOnHand:       input.QtyOnHand,
		BrandName:       input.BrandName,
		IsActive:        input.IsActive,
	}

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

// GetProducts godoc
//
//	@Summary		Fetch all products
//	@Description	Fetch all products
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Success		200				{array}		models.Product
//	@Failure		400				{object}	httputil.HttpError400
//	@Failure		401				{object}	httputil.HttpError401
//	@Failure		500				{object}	httputil.HttpError500
//	@Router			/api/product	[get]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
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

// GetProductById godoc
//
//	@Summary		Fetch individual product by Id
//	@Description	Fetch individual product by Id
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string	true	"product Id"
//	@Success		200					{object}	models.Product
//	@Failure		400					{object}	httputil.HttpError400
//	@Failure		401					{object}	httputil.HttpError401
//	@Failure		500					{object}	httputil.HttpError500
//	@Router			/api/product/{id}	[get]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func GetProductById(c *fiber.Ctx) error {
	db := database.DB

	id := strings.ToUpper(c.Params("id"))

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

// UpdateProduct godoc
//
//	@Summary		Update individual product
//	@Description	Update individual product
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string						true	"product Id"
//	@Param			product				body		dto.UpdateProductRequstDTO	true	"Product Data"
//	@Success		200					{object}	models.Product
//	@Failure		400					{object}	httputil.HttpError400
//	@Failure		401					{object}	httputil.HttpError401
//	@Failure		500					{object}	httputil.HttpError500
//	@Router			/api/product/{id}	[put]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
func UpdateProduct(c *fiber.Ctx) error {

	db := database.DB

	input := new(dto.UpdateProductRequstDTO)
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
		// product.QtyOnHand = input.QtyOnHand
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

// DeleteProduct godoc
//
//	@Summary		Delete individual product
//	@Description	Delete individual product
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string	true	"product Id"
//	@Success		200					{object}	models.Product
//	@Failure		400					{object}	httputil.HttpError400
//	@Failure		401					{object}	httputil.HttpError401
//	@Failure		500					{object}	httputil.HttpError500
//	@Router			/api/product/{id}	[delete]
//
//	@Security		ApiKeyAuth
//
//	@Security		Bearer  <-----------------------------------------add this in all controllers that need authentication
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
