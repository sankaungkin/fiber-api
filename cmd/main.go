package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sankaungkin/fiber-api/database"
	"github.com/sankaungkin/fiber-api/router"
)



func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	//create fiber app
	app := fiber.New()

	//connect database
	database.ConnectDB()

	//create routes
	router.Initialize(app)
	app.Listen(":5555")
}