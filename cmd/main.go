package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger" // swagger handler
	"github.com/joho/godotenv"

	_ "github.com/sankaungkin/fiber-api/cmd/docs"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/router"
)

//	@title			Fiber API
//	@version		1.0
//	@description	This is a sample swagger for Fiber
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	//create fiber app
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	//connect database
	database.ConnectDB()

	//create routes
	router.Initialize(app)
	app.Listen(":5555")
}
