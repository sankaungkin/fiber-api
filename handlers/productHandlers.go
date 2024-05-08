package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	_ "github.com/sankaungkin/fiber-api/cmd/docs"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/gorm"
)

// CreateProduct 	godoc
//
//	@Summary		Create new product based on paramters
//	@Description	Create new product based on paramters
//	@Tags			Products
//	@Accept			json
//	@Param			product	body		models.Product	true	"Product Data"
//	@Success		200		{object}	models.Product
//	@Failure		400,500	{object}	httputil.HttpError
//	@Router			/api/product [post]
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

// GetProducts godoc
//
//	@Summary		Fetch all products
//	@Description	Fetch all products
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Success		200				{array}		models.Product
//	@Failure		400,500			{object}	httputil.HttpError
//	@Router			/api/product	[get]
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
//	@Failure		400,500				{object}	httputil.HttpError
//	@Router			/api/product/{id}	[get]
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

// UpdateProduct godoc
//
//	@Summary		Update individual product
//	@Description	Update individual product
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string					true	"product Id"
//	@Param			product				body		models.UpdateProductDTO	true	"Product Data"
//	@Success		200					{object}	models.Product
//	@Failure		400,500				{object}	httputil.HttpError
//	@Router			/api/product/{id}	[put]
func UpdateProduct(c *fiber.Ctx) error {

	db := database.DB

	input := new(models.UpdateProductDTO)
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
//	@Summary		Update individual product
//	@Description	Update individual product
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string	true	"product Id"
//	@Success		200					{object}	models.Product
//	@Failure		400,500				{object}	httputil.HttpError
//	@Router			/api/product/{id}	[delete]
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
