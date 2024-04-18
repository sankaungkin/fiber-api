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
	categories.Post("/", middleware.Protected(), handlers.CreateCategory)
	categories.Get("/", middleware.Protected(), handlers.GetCategories)
	categories.Get("/:id", middleware.Protected(), handlers.GetCategory)

	// Auth
	users := api.Group("/user")
	users.Post("/", handlers.CreateUser)
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Post("/login", handlers.Login)
}
