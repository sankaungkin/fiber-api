package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sankaungkin/fiber-api/models"
	"github.com/sankaungkin/fiber-api/storage"
	"gorm.io/gorm"
)

type Repository struct{
	DB *gorm.DB
}

type Category struct {
	
	CategoryName	string	`json:"categoryName"`

}

func(r *Repository) CreateCategory(c *fiber.Ctx) error {
	category := Category{}

	// validate the CREATE CATEGORY DTO
	err := c.BodyParser(&category)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{ "message":"request failed"})
		return err
	}

	err = r.DB.Create(&category).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{ "message": "could not create new category"})
		return err
	}
	 
	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"category has been created successfully"})
	return nil
	}

func(r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/category", r.CreateCategory)
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from FIBER API -> <---")
	})
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("API_PORT"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		User: os.Getenv("POSTGRES_USER"),
		DBName: os.Getenv("POSTGRES_DB"),
		SSLMode: os.Getenv("SSLMODE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not get connection to database")
	}

	err = models.MigrateCategory(db)
	if err != nil {
		log.Fatal("could not migrate category table to db")
	}


	r := Repository {
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":5555")
}