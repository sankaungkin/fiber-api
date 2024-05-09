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

// @title						Fiber API
// @version					1.0
// @description				This is an auto-generated API docs.
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.email				sankaungkin@gmail.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
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
