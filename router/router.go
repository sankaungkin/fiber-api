package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/sankaungkin/fiber-api/handlers"
)

func Initialize(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello from stt api using go fiber framework.. <-- ")
	})

	categories := app.Group("/category")
	categories.Post("/", handlers.CreateCategory)
	categories.Get("/", handlers.GetCategories)
	categories.Get("/:id", handlers.GetCategory)

}
