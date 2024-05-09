package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/sankaungkin/fiber-api/handlers"
	"github.com/sankaungkin/fiber-api/middleware"
)

func Initialize(app *fiber.App) {

	api := app.Group("/api", logger.New())
	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello from stt api using go fiber framework.. <-- ")
	})

	// category
	categories := api.Group("/category")
	categories.Use(middleware.Protected())
	categories.Post("/", handlers.CreateCategory)
	categories.Get("/", handlers.GetCategories)
	categories.Get("/:id", handlers.GetCategoryById)
	categories.Put("/:id", handlers.UpdateCategory)
	categories.Delete("/:id", handlers.DeleteCategory)

	// Auth
	users := api.Group("/auth")
	users.Post("/signup", handlers.CreateUser)
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Post("/login", handlers.Login)
	users.Post("/logout", handlers.Logout)
	users.Post("/refresh", handlers.Refresh)

	// product
	products := api.Group("/product")
	products.Use(middleware.Protected())
	products.Get("/", handlers.GetProducts)
	products.Get("/:id", handlers.GetProductById)
	products.Post("/", middleware.Authorize, handlers.CreateProduct)
	products.Put("/:id", middleware.Authorize, handlers.UpdateProduct)
	products.Delete("/:id", middleware.Authorize, handlers.DeleteProduct)

	// Inventory
	inventories := api.Group("/inventory")
	// inventories.Use(middleware.Protected())
	inventories.Post("/increase", middleware.Authorize, handlers.IncreaseInventory)
	inventories.Post("/decrease", middleware.Authorize, handlers.DecreaseInventory)

	// Sales
	sales := api.Group("/sale")
	sales.Post("/", handlers.CreateSaleInvoice)
	sales.Get("/", handlers.GetSales)

	// Purchases
	purchases := api.Group("/purchase")
	purchases.Post("/", handlers.CreatePurchaseInvoice)
	purchases.Get("/", handlers.GetPurchases)

}
